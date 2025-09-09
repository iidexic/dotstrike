package pops

import (
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"iidexic.dotstrike/config"
)

// // TODO: Document this (i dont remember what it's for)
// type esource = byte
//
// const (
//
//	_ esource = iota
//	esInfileOPEN
//	esOutfileMAKEOPEN
//	esCOPY
//	es
//
// )
type DirEntry = fs.DirEntry
type boolConfig map[config.OptionKey]bool

func (b boolConfig) IsOn(k config.OptionKey) bool { v, ok := b[k]; return ok && v }

func (b *boolConfig) CopyMap(m map[config.OptionKey]bool) {
	maps.Copy(*b, m)
}

type stringConfig = map[config.OptionKey]string
type PathError = fs.PathError

var (
	globalOutDir *string //NOTE: Using  copierMaschine GlobalOut instead
	bNoFiles     = config.BoolNoFiles
	bAllDirs     = config.BoolCopyAllDirs
	bRootSubdir  = config.BoolRootSubdir
	bUseGlobal   = config.BoolUseGlobalTarget
	bNoHidden    = config.BoolIgnoreHidden
	bNoRepo      = config.BoolIgnoreRepo
)

// copierMaschine builds and executes CopyJobs
// Ideally single-instance; use GetCopier to get
type copierMaschine struct {
	JobQueue     map[string]*CopyJob
	JobGroups    map[string]*JobGroup
	globalOut    string
	setglobalOut bool
}

var cmachine copierMaschine = copierMaschine{
	JobQueue:  make(map[string]*CopyJob),
	JobGroups: make(map[string]*JobGroup),
}

func Copier() *copierMaschine { return &cmachine }

type JobGroup struct {
	pathSet
	groupName   string
	jobNames    []string
	jobPtrs     []*CopyJob
	bcfg        boolConfig
	scfg        stringConfig
	initialized bool
}

type pathSet struct {
	ins  []string
	outs []string
}

// type jobConfig struct {
// 	bcfg boolConfig
// 	scfg stringConfig
// }

// TODO: finish RunAll

func (CM *copierMaschine) RunAll(stopOnError bool) {
	for name, job := range CM.JobQueue {
		_ = job
		_ = name
	}
}

func (CM *copier) SetGlobalOutDir(globalOut string) {

}

func (CM *copierMaschine) Detail() []string {
	d := make([]string, len(CM.JobQueue)+2)
	return d
}

func (g *JobGroup) Detail() string {
	sd := make([]string, len(g.jobPtrs)+len(g.bcfg)+len(g.scfg)+2)
	sd[0] = fmt.Sprintf("-[Group: %s] %d jobs", g.groupName, len(g.jobPtrs))
	i := 1
	if len(g.jobPtrs) > 0 {
		sd[i] = "-- JOBS --"
		i++
		for _, j := range g.jobPtrs {
			sd[i+1] = fmt.Sprintf("---[%d] ", i) + j.DetailLine() + "\n"
			i++
		}
	}
	if len(g.scfg)+len(g.bcfg) > 0 {
		sd[i] = "-- GROUP CONFIG --"
		i++
		for k, bv := range g.bcfg {
			sd[i] = fmt.Sprintf("--- %s: %t", k.String(), bv)
			i++
		}
		for k, v := range g.scfg {
			sd[i] = fmt.Sprintf("--- %s: '%s'", k.String(), v)
		}
	} else {
		sd = slices.Clip(sd)
	}

	return strings.Join(sd, "\n")
}

// TODO:(low) finish full CopyJob Detail
func (J *CopyJob) Detail() []string {
	d := make([]string, 4+len(J.fstack)+len(J.ignore.Patterns), +len(J.BPrefs)+len(J.BPrefs))
	d[0] = fmt.Sprintf("in:'%s' out: %s | ", J.PathIn, J.PathOut)

	return d
}

func (J *CopyJob) DetailLine() string {
	d := fmt.Sprintf("in:'%s' | out: %s ", J.PathIn, J.PathOut)
	if len(J.ignore.Patterns) > 0 {
		d += fmt.Sprintf("| #ignores:%d", len(J.ignore.Patterns))
	}
	if J.jobRan {
		d += fmt.Sprintf("| ran (%d file, %d newdir, %d errors)", len(J.fstack), J.DirsMade(), len(J.OpErrors))
	}
	return d
}

// NewJobGroup takes all required data for a group of related Copy Jobs, and returns a JobGroup ptr.
//
// It also automatically creates all copy jobs, and stores the job names in JobGroup.jobNames.
// Job names are created as (name-[job#])

func (CM *copierMaschine) NewJobGroup(name string, inPaths []string, outPaths []string, bools boolConfig) *JobGroup {
	numJobs := len(inPaths) * len(outPaths)
	group := &JobGroup{groupName: name, bcfg: bools,
		jobNames: make([]string, numJobs),
		jobPtrs:  make([]*CopyJob, numJobs),
		pathSet:  pathSet{ins: inPaths, outs: outPaths},
	}
	//TODO: Must process UseGlobal+KillGlobal now
	group.makeJobs()
	group.initialized = true
	key := group.groupName
	_, ok := CM.JobGroups[key]
	if ok {
		n := 1
		for {
			key += fmt.Sprintf("%2d", n)
		}
	}
	CM.JobGroups[group.groupName] = group
	return group
}

func (g *JobGroup) makeJobs() {
	for x, in := range g.ins {
		for y, out := range g.outs {
			i := x*len(g.outs) + y
			jname := fmt.Sprintf("%s.in-%d.out-%d", g.groupName, x, y)
			g.jobPtrs[i] = cmachine.NewJob(jname, in, out)
			g.jobNames[i] = jname
		}
	}

}

func (g *JobGroup) RunAll(abortOnError bool) error {
	var outError error
	for i := range g.jobPtrs {
		e := g.jobPtrs[i].Run()
		if e != nil && abortOnError {
			return e
		} else if e != nil {
			if outError == nil {
				outError = fmt.Errorf("Copy Errors: %w", e)
			} else {
				outError = fmt.Errorf("%w\n%w", outError, e)
			}
		}
	}
	return outError
}

// CopyJob prepares and executes the copy of all contents of PathIn to PathOut
type CopyJob struct {
	PathIn, PathOut string     // Root of copy source and destination (*or destination parent)
	parentPathOut   string     // unused. populated on run if JobSettings.makeRootSubdir = true
	fstack          []filecopy // record of files copied
	jobRan          bool
	newDirs         map[string]bool //
	ignore          IgnoreSet
	OpErrors        []fs.PathError
	BPrefs          boolConfig
	SPrefs          stringConfig //* Currently Unused
}

// NewJob creates a new job and adds to the JobQueue; returns a ptr to created job if successful
// jobName must be unique within the JobQueue; if NewJob is passed an existing jobName,
// it will not add the job to the queue, and will return nil
func (CM *copierMaschine) NewJob(jobName, pathIn, pathOut string) *CopyJob {
	if CM.jobExists(jobName) {
		return nil
	}
	CM.JobQueue[jobName] = &CopyJob{PathIn: pathIn, PathOut: pathOut, newDirs: make(map[string]bool)}
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

func (CM *copierMaschine) groupExists(jobName string) bool {
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

func (J *CopyJob) configCheck(opt config.OptionKey) bool {
	if opt.IsBool() {
		v, ok := J.BPrefs[opt]
		return v && ok
	}
	if opt.IsString() {
		v, ok := J.SPrefs[opt]
		return len(v) > 0 && ok
	}

	return false
}

// TODO:  test IgnoreGit  (ALSO Test Global)
func (J *CopyJob) IgnoreGit() {
	J.ignore.Patterns = append(J.ignore.Patterns, subptn{ptn: `.git`, matchDir: true})
}

// filecopy acts as a record of a single file's copy operation
type filecopy struct {
	relpath         string
	inSize, outSize int64
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

// ── Performing CopyJob ──────────────────────────────────────────────

// Run the copy. Returns
func (J *CopyJob) Run( /* params */ ) error {
	// condition paths
	var e error
	// Prefs that need to be processed here
	//TODO: CHECK WHETHER OR NOT BPREFS ALWAYS CONTAINS ALL OF THESE VALUES
	//	IF NOT MUST ADD CHECK FOR KEY EXIST TO EVERYTHING
	/*
	   bNoFiles
	   bAllDirs
	   bRootSubdir
	   bUseGlobal
	   bNoHidden
	   bNoRepo
	*/
	if J.BPrefs.IsOn(bNoRepo) {

	}
	if J.BPrefs.IsOn(bNoHidden) {

	}

	if J.jobRan {
		return fmt.Errorf("CopyJob Already Ran")
	}
	if J.BPrefs[bRootSubdir] && J.parentPathOut == "" { // parentpath check to be safe
		J.parentPathOut = J.PathOut
		J.PathOut = Joinpath(J.PathOut, filepath.Base(J.PathIn))
	}
	J.jobRan = true
	// ── Walk ────────────────────────────────────────────────────────────
	e = filepath.WalkDir(J.PathIn, J.Walk)

	if e != nil {
		J.OpErrors = append(J.OpErrors, fs.PathError{Path: J.PathIn, Err: e, Op: ""})
		return e
	}

	//warning: WITH THIS STRUCTURE, A WALKDIR ERROR WILL PREVENT MAKING ADDITIONAL DIRS
	// why would I do this separately?
	// Technically it works.. but I don't see a reason not to have this in the walk
	if J.BPrefs[bAllDirs] {
		for relDir := range J.newDirs {
			e := os.MkdirAll(Joinpath(J.PathOut, relDir), 0)
			J.checkAndLogError(relDir, "MakeDirectory", e)
		}
	}
	return nil
}

// logDir adds directories to j.newDirs if they are not already present
// NOTE: Walk sends relative paths to logDir (J.newDirs keys will be relative)
func (J *CopyJob) logDir(dir string, copied bool) {
	var exists bool
	for keydir := range J.newDirs {
		exists = (exists || keydir == dir)
	}
	if !exists {
		J.newDirs[dir] = copied
	}
}
func (J *CopyJob) DirsMade() int {
	n := 0
	for _, v := range J.newDirs {
		if v {
			n++
		}
	}
	return n
}

// ╭─────────────────────────────────────────────────────────╮
// │                      WALK FUNCTION                      │
// ╰─────────────────────────────────────────────────────────╯

func (J *CopyJob) Walk(p string, d DirEntry, e error) error {
	// make relative path first; used for dirs & files
	rootRelativePath := J.stripRoot(p) //	!INFO: panics on error; error is unexpected

	// DIRECTORIES:
	if d.IsDir() {
		// check ignore + prevent recursion (if PathOut is a subdir of PathIn)
		if J.ignore.isIgnored(p, true) || strings.HasPrefix(p, J.PathOut) {
			J.logDir(rootRelativePath, false)
			return fs.SkipDir
		}
		J.logDir(rootRelativePath, true)
		return nil
	} else { // File Ignore
		if J.ignore.isIgnored(p, false) {
			return nil
		}
	}

	// ── 0.1 make filepath out ─────────────────────────────────
	pto := Joinpath(J.PathOut, rootRelativePath)
	if !filepath.IsAbs(pto) {
		pto = absNoE(pto) //	!INFO: panics on error; error is unexpected
	}

	// 0.2 Get infile Info()
	inDE, e := d.Info()
	J.checkAndLogError(p, "GetFileInfo_In", e)

	// 0.3 Return before copy if config requires.
	if J.BPrefs[bNoFiles] {
		J.addFile(rootRelativePath, inDE.Size(), 0) // for dry runs
		return nil
	}

	// ── 1. open in file ──
	inF, e := OpenExistingFile(p)
	defer inF.Close()
	if e != nil {
		J.logError(p, "OpenExisting_In", e)
		return nil //skip file if opening errors
	}
	// ── 2. make/open out file ──

	outF, e := MakeOpenFileF(pto)
	defer outF.Close()
	J.checkAndLogError(pto, "MakeOpen_Out", e)

	// ── 3. perform copy ──
	wb, e := io.CopyBuffer(outF, inF, nil)
	J.checkAndLogError(pto, "CopyError_Out", e)

	// check size original matches new copied
	J.addFile(rootRelativePath, inDE.Size(), wb)
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
