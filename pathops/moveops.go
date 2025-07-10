package pops

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
)

type esource = byte

const (
	_ esource = iota
	esInfileOPEN
	esOutfileMAKEOPEN
	esCOPY
	es
)

type DirEntry = fs.DirEntry

// copierMaschine builds and executes CopyJobs
// Ideally single-instance; use GetCopier to get
type copierMaschine struct {
	JobQueue map[string]*CopyJob
}

// primary-instance use
var cmachine copierMaschine = copierMaschine{JobQueue: make(map[string]*CopyJob)}

func GetCopierMaschine() *copierMaschine { return &cmachine }

// NewJob creates a new job and adds to the JobQueue; returns true if successful
// jobName must be unique within the JobQueue; if NewJob is passed an existing jobName,
// it will not add the job to the queue, and will return false
func (CM *copierMaschine) NewJob(jobName, pathIn, pathOut string) *CopyJob {
	if CM.jobExists(jobName) {
		return nil
	}
	CM.JobQueue[jobName] = &CopyJob{PathIn: pathIn, PathOut: pathOut}
	return CM.JobQueue[jobName]
}

func (CM *copierMaschine) jobExists(jobName string) bool {
	for k := range CM.JobQueue {
		if k == jobName {
			return true
		}
	}
	return false
}

// GetJob returns *CopyJob if jobName exists in the JobQueue
// otherwise returns nil ptr
func (CM *copierMaschine) GetJob(jobName string) *CopyJob {
	for keyName, ptrjob := range CM.JobQueue {
		if keyName == jobName {
			return ptrjob
		}
	}
	return nil
}

// RunJob is equivalent to running GetJob and then running CopyJob.Run(nil)
func (CM *copierMaschine) RunJob(jobName string) *CopyJob {
	if ptrjob := CM.GetJob(jobName); ptrjob != nil {
		ptrjob.Run()
		return ptrjob
	}
	return nil
}

// CopyJob prepares and executes the copy of all contents of PathIn to PathOut
type CopyJob struct {
	PathIn, PathOut string // Root of copy source and destination (*or destination parent)
	originalPathOut string // unused. populated on run if JobSettings.makeRootSubdir = true
	fstack          []filecopy
	ignore          IgnoreSet
	OpErrors        []fs.PathError
	JobSettings     copyConfig
}

type copyConfig struct {
	makeRootSubdir bool // if true, appends base(PathIn) to PathOut
	//DRY_RUN bool
}

// SetOptionMakeSubdir - sets CopyJob.JobSettings.makeRootSubdir
// this appends the PathIn dir name to PathOut
// if true: set PathOut= filepath.Join(PathOut, filepath.Base(PathIn)), store original
// if false: Copy directly into provided PathOut
func (J *CopyJob) SetOptionMakeSubdir(makeSubdir bool) {
	J.JobSettings.makeRootSubdir = makeSubdir
}

// filecopy acts as a record of a single file's copy operation
type filecopy struct {
	relpath         string
	inSize, outSize int64
}

// joinpath wraps filepath.Join, but fixes os path inconsistency before returning
func joinpath(elem ...string) string {
	p := filepath.Join(elem...)
	//NOTE: Commenting out for testing
	/* if runtime.GOOS == "windows" {
		p = strings.Replace(p, "/", `\`, -1)
	} */
	return p
}

// Removed for now.
// type elog struct {
// 	fs.PathError
// 	relpath string
// 	esrc    esource
// }

func (J *CopyJob) logError(abspath, opname string, e error) {
	J.OpErrors = append(J.OpErrors, fs.PathError{Path: abspath, Op: opname, Err: e})
}

// checkAndLogError checks the error, and logs non-nil errors to CopyJob.logError.
// returns true if error!=nil, else false
func (J *CopyJob) checkAndLogError(abspath, opname string, e error) bool {
	if e != nil {
		J.logError(abspath, opname, e)
		return true
	}
	return false
}

func (J *CopyJob) addFile(relpath string, inSize, outSize int64) {
	J.fstack = append(J.fstack, filecopy{relpath: relpath, inSize: inSize, outSize: outSize})
}

func (J *CopyJob) Run( /* params */ ) error {
	// condition paths
	var e error
	//NOTE: Commenting out for testing
	/* J.PathIn, e = conditionPath(J.PathIn)
	if J.checkAndLogError(J.PathIn, "conditionPath", e) {
		return fmt.Errorf("error abs(PathIn): %w", e) }
	J.PathOut, e = conditionPath(J.PathOut)
	if J.checkAndLogError(J.PathOut, "conditionPath", e) {
		return fmt.Errorf("error abs(PathOut): %w", e)
	} */

	// Add folder name to PathOut if needed. *Run this after clean
	// checking for blank originalPathOut just to avoid potential future issue
	if J.JobSettings.makeRootSubdir && J.originalPathOut == "" {
		J.originalPathOut = J.PathOut
		J.PathOut = joinpath(J.PathOut, filepath.Base(J.PathIn))
		// Join uses "OS specific Separator"
	}
	e = filepath.WalkDir(J.PathIn, J.Walk)
	if e != nil {
		J.OpErrors = append(J.OpErrors, fs.PathError{Path: J.PathIn, Err: e, Op: ""})
		return e
	}
	return nil
}

func (J *CopyJob) Walk(p string, d DirEntry, e error) error {
	// removed filepath.clean(p)

	// DIRECTORIES:
	if d.IsDir() {
		// check ignore + prevent recursion (if PathOut is a subdir of PathIn)
		if J.ignore.isIgnored(p) || strings.HasPrefix(p, J.PathOut) {
			return fs.SkipDir
		}
		return nil
	}
	// FILES:
	// ── 0. make filepath out ─────────────────────────────────
	prr := J.stripRoot(p) //	!INFO: panics on error; error is unexpected(expand)
	pto := joinpath(J.PathOut, prr)
	if !filepath.IsAbs(pto) {
		pto = absNoE(pto) //	!INFO: panics on error; error is unexpected(expand)
	}
	// ── 1. open in file ──────────────────────────────────────
	inF, e := OpenExistingFile(p)
	defer inF.Close()
	if e != nil {
		J.logError(p, "OpenExisting_In", e)
		return nil //skip file if opening errors
	}
	// ── 2. make+open out file ────────────────────────────────
	outF, e := MakeOpenFileF(pto)
	defer outF.Close()
	J.checkAndLogError(pto, "MakeOpen_Out", e)
	// ── 3. perform copy ──────────────────────────────────────
	wb, e := io.CopyBuffer(outF, inF, nil)
	J.checkAndLogError(pto, "CopyError_Out", e)

	// check size original matches new copied
	inDE, e := d.Info()
	J.checkAndLogError(p, "GetFileInfo_In", e)
	J.addFile(prr, inDE.Size(), wb)
	return nil
}

// stripRoot removes CopyJob.PathIn from path provided for construction of destination path
// structure/intent of CopyJob requires J.PathIn to be a prefix in rpath.
// As such, if an error is encountered, stripRoot panics
func (J *CopyJob) stripRoot(p string) string {
	relp, e := filepath.Rel(J.PathIn, p)
	if e != nil {
		panic(fmt.Errorf("stripRoot(%s) error: %v", p, e))
	}
	return relp

}

// absNoE runs abs and returns the resulting string; panics on error
func absNoE(p string) string {
	po, e := filepath.Abs(p)
	if e != nil {
		panic(e)
	}
	return po
}

// conditionPath cleans path, gets abs path, and fixes windows path separators
// returns error direct from filepath.Abs
func conditionPath(p string) (string, error) {
	var e error
	p = filepath.Clean(p)
	p, e = filepath.Abs(p)
	if runtime.GOOS == "windows" {
		p = strings.Replace(p, "/", `\`, -1)
	}
	return p, e
}

// ── Ignore Functionality ────────────────────────────────────────────
// TODO: Finish Ignore system. For now it's not priority
// IgnoreSet stores & processes ignore data for a CopyJob
type IgnoreSet struct {
	ignoreDir  []iptn
	ignoreFile []iptn
}

// iptn is a single ignore string pattern
// pattern is loaded with raw provided string
// anyL + anyR are true if pattern[0]=="*" and pattern[len-1]=="*"  respectively
type iptn struct {
	pattern    string
	anyL, anyR bool
	dir        bool
	psize      byte
}

// matches checks string against the valid iptn
func (ip iptn) matches(s string) bool {
	if strings.Index(s, ip.pattern) >= 0 {
		return true
	}

	return false
}

func (I *IgnoreSet) isIgnored(dirpath string) bool {
	for _, pat := range I.ignoreDir {
		if pat.matches(dirpath) {
			return true
		}
	}
	return false
}
