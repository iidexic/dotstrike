package pops

import (
	"fmt"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"iidexic.dotstrike/config"
)

type DirEntry = fs.DirEntry
type boolConfig map[config.OptionKey]bool

func (b boolConfig) IsOn(k config.OptionKey) bool { v, ok := b[k]; return ok && v }

func (b *boolConfig) CopyMap(m map[config.OptionKey]bool) {
	maps.Copy(*b, m)
}

type stringConfig = map[config.OptionKey]string
type PathError = fs.PathError

var (
	bNoFiles    = config.BoolNoFiles
	bAllDirs    = config.BoolCopyAllDirs
	bRootSubdir = config.BoolRootSubdir
	bUseGlobal  = config.BoolUseGlobalTarget
	bNoHidden   = config.BoolIgnoreHidden
	bNoRepo     = config.BoolIgnoreRepo
)

// copierMaschine builds and executes CopyJobs
// Ideally single-instance; use GetCopier to get
type copierMaschine struct {
	JobQueue  map[string]*CopyJob
	JobGroups map[string]*JobGroup
	runErrors []error
}

var cmachine copierMaschine = copierMaschine{
	JobQueue:  make(map[string]*CopyJob),
	JobGroups: make(map[string]*JobGroup),
	runErrors: make([]error, 0, 32),
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

func (CM *copierMaschine) RunAll(stopOnError bool) error {
	for name, group := range CM.JobGroups {
		e := group.RunAll(stopOnError)
		if e != nil {
			if stopOnError {
				return fmt.Errorf("error for %s: %w", name, e)
			} else {
				CM.runErrors = append(CM.runErrors, e)

			}
		}
	}
	//TODO: Return Error?
	return nil
}

func (CM *copierMaschine) Detail() []string {
	dtlength := max(len(CM.JobQueue), len(CM.JobGroups)) + 2
	if dtlength == 2 {
		dtlength = 6 //safety
	}
	d := make([]string, dtlength)
	d[0] = "|[Copier State]---------------------|"
	n := 1
	for _, g := range CM.JobGroups {
		if n < len(d) {
			d[n] = g.Detail()
		} else {
			d = append(d, g.Detail())
		}
		n++
	}
	d[n] = "|-----------------------------------|"
	d = slices.Clip(d)
	return d
}

func (g *JobGroup) Detail() string {
	sd := make([]string, len(g.jobPtrs)+len(g.bcfg)+len(g.scfg)+3)
	sd[0] = fmt.Sprintf("|[Group: %s] %d jobs", g.groupName, len(g.jobPtrs))
	i := 1
	if len(g.jobPtrs) > 0 {
		sd[i] = "|--- JOBS:"
		i++
		for ijp, j := range g.jobPtrs {
			detailText := fmt.Sprintf("|---	[%d] ", ijp) + j.DetailLine()
			sd[i] = detailText
			i++
		}
	}
	if len(g.scfg)+len(g.bcfg) > 0 {
		sd[i] = "|-- GROUP CONFIG:"
		i++
		for k, bv := range g.bcfg {
			sd[i] = fmt.Sprintf("|---	%s: %t", k.String(), bv)
			i++
		}
		for k, v := range g.scfg {
			sd[i] = fmt.Sprintf("|---	%s: '%s'", k.String(), v)
		}
	}
	sd = slices.Clip(sd)
	return strings.Join(sd, "\n")
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
		e := g.jobPtrs[i].RunFS()
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

// NewJob creates a new job and adds to the JobQueue; returns a ptr to created job if successful
// jobName must be unique within the JobQueue; if NewJob is passed an existing jobName,
// it will not add the job to the queue, and will return nil
func (CM *copierMaschine) NewJob(jobName, pathIn, pathOut string) *CopyJob {
	if CM.jobExists(jobName) {
		return nil
	}
	CM.JobQueue[jobName] = &CopyJob{PathIn: pathIn, PathOut: pathOut,
		newDirs: make(map[string]bool),
		BPrefs:  make(boolConfig),
		SPrefs:  make(stringConfig),
		record:  fileRecord{files: make(map[string]*filedata)}}
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
		ptrjob.RunFS()
		return ptrjob
	}
	return nil
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

func (F filecopy) String() string {
	var is, os float64
	var iu, ou string
	if F.inSize > 1048576 {
		is = float64(F.inSize) / 1048576.0
		iu = "MB"
	} else {
		is = float64(F.inSize) / 1024.0
		iu = "KB"
	}
	if F.outSize > 1048576 {
		os = float64(F.outSize) / 1048578.0
		ou = "MB"
	} else {
		os = float64(F.outSize) / 1024.0
		ou = "KB"
	}
	return fmt.Sprintf("'%s' (In: %.2f %s, Out: %.2f %s)", F.relpath, is, iu, os, ou)
}

// absNoE runs abs and returns the resulting string; panics on error
func absNoE(p string) string {
	po, e := filepath.Abs(p)
	if e != nil {
		panic(e)
	}
	return po
}

// TODO: Remove conditionPath

// conditionPath cleans path and gets abs path
// returns error direct from filepath.Abs
func conditionPath(p string) (string, error) {
	var e error
	p = filepath.Clean(p)
	p, e = filepath.Abs(p)
	// if runtime.GOOS == "windows" {
	// 	p = strings.Replace(p, "/", `\`, -1)
	// }
	return p, e
}
