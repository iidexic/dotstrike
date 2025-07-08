package pops

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
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
var cmachine copierMaschine

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
// Pre-Run:
// PathIn & PathOut are root dir of original files and copy location, respectively
// ignore contains list of files to ignore; must be populated before running
// JobSettings contains job-specific settings that will affect how copy job runs
// Post-Run:
// fstack will be populated when the job is run, providing a record of what is copied
// OpErrors contains encountered errors, using fs.PathError format.
type CopyJob struct {
	PathIn, PathOut string
	originalPathOut string // will be populated if makeRootSubdir = true
	fstack          []filecopy
	ignore          IgnoreSet
	OpErrors        []fs.PathError
	JobSettings     copyConfig
}

type copyConfig struct {
	makeRootSubdir bool
	//DRY_RUN bool
}

// SetOptionMakeSubdir - sets CopyJob.JobSettings.makeRootSubdir
// this appends the PathIn dir name to PathOut
// if true: given PathIn == (base-pathIn...\DirName), set PathOut = (PathOut\DirName)
// if false: Copy directly to provided PathOut
func (J *CopyJob) SetOptionMakeSubdir(makeSubdir bool) {
	J.JobSettings.makeRootSubdir = makeSubdir
}

// filecopy acts as a record of a single file's copy operation
type filecopy struct {
	name, relpath   string
	id              int //?
	inPath, outPath string
	inSize, outSize int64
}

// Removed for now.
// type elog struct {
// 	fs.PathError
// 	relpath string
// 	esrc    esource
// }

func (J *CopyJob) Load(absPathIn, absPathOut string) {
	J.fstack = append(J.fstack, filecopy{
		inPath: absPathIn, outPath: absPathOut,
	})
}
func (J *CopyJob) Run( /* params */ ) error {
	if J.JobSettings.makeRootSubdir {
		J.originalPathOut = J.PathOut
		J.PathOut = filepath.Join(J.PathOut, filepath.Base(J.PathIn))
	}
	e := filepath.WalkDir(J.PathIn, J.Walk)
	if e != nil {
		J.OpErrors = append(J.OpErrors, fs.PathError{Path: J.PathIn, Err: e, Op: ""})
		return e
	}
	return nil
}

func (J *CopyJob) logError(abspath, opname string, e error) {
	J.OpErrors = append(J.OpErrors, fs.PathError{Path: abspath, Op: opname, Err: e})
}

func (J *CopyJob) checkAndLogError(abspath, opname string, e error) {
	if e != nil {
		J.logError(abspath, opname, e)
	}
}

// NOTE: Removed. Amount of detail unnecessary
func (J *CopyJob) newFilecopy(inInfo *fs.FileInfo, inPath, outPath string, outSize int64) *filecopy {
	J.fstack = append(J.fstack, filecopy{name: (*inInfo).Name(),
		inPath: inPath, outPath: outPath, outSize: outSize})
	return &J.fstack[len(J.fstack)-1]
}

func (J *CopyJob) addFile(relpath string, inSize, outSize int64) {
	J.fstack = append(J.fstack, filecopy{relpath: relpath, inSize: inSize, outSize: outSize})
}

func (J *CopyJob) Walk(p string, d DirEntry, e error) error {
	// DIRECTORIES:
	if d.IsDir() {
		if J.ignore.isIgnored(p) {
			return fs.SkipDir
		}
		return nil
	}
	// FILES:
	// ── 0. make filepath out ─────────────────────────────────
	prel := J.stripRoot(p) //	!INFO: panics on error; error is unexpected(expand)
	prr := prel
	// add folder name if option enabled
	if J.JobSettings.makeRootSubdir {
		prr = filepath.Join(filepath.Base(J.PathIn), prr)
	}
	pto := filepath.Join(J.PathOut, prr)
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
	J.addFile(prel, inDE.Size(), wb)
	return nil
}

// absNoE runs abs and returns the resulting string; panics on error
func absNoE(p string) string {
	po, e := filepath.Abs(p)
	if e != nil {
		panic(e)
	}
	return po
}

// stripRoot removes CopyJob.PathIn from arg path `rpath` for construction of destination path
// structure/intent of CopyJob requires J.PathIn to be a prefix in rpath.
// As such, if an error is encountered, stripRoot panics
func (J *CopyJob) stripRoot(rpath string) string {
	p, e := filepath.Rel(rpath, J.PathIn)
	if e != nil {
		panic(fmt.Errorf("stripRoot(%s) error: %v", rpath, e))
	}
	return p

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
