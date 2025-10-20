package pops

import (
	"fmt"
	"io/fs"
	"strconv"
	"strings"

	"iidexic.dotstrike/config"
	"iidexic.dotstrike/uout"
)

type DirEntry = fs.DirEntry
type boolConfig = config.ConfigMap

//func (b boolConfig) IsOn(k config.OptionKey) bool { v, ok := b[k]; return ok && v }

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
func (CM *copierMaschine) String() string { return strings.Join(CM.Detail(), "\n") }

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
	d = d[:n+1]
	return d
}

func (CM copierMaschine) GroupDetails() string {
	out := uout.NewOut("--[Copier Job Groups]--")
	out.IndR()
	out.ILV(CM.JobGroups)

	return out.String()
}

// NewJobGroup takes all required data for a group of related Copy Jobs, and returns a JobGroup ptr.
//
// It also automatically creates all copy jobs, and stores the job names in JobGroup.jobNames.
// Job names are created as (name-[job#])
func (CM *copierMaschine) NewJobGroup(UniqueName string, inPaths []string, outPaths []string, bools boolConfig) *JobGroup {
	numJobs := len(inPaths) * len(outPaths)
	group := &JobGroup{groupName: UniqueName, bcfg: bools,
		jobNames: make([]string, numJobs),
		jobPtrs:  make([]*CopyJob, numJobs),
		pathSet:  pathSet{ins: inPaths, outs: outPaths},
	}

	group.makeJobs()
	group.initialized = true
	//  what is happening here
	key := group.groupName
	_, ok := CM.JobGroups[key]

	// lazy way to make sure group names are unique.
	n := 0
	k := key
	for {
		if ok {
			n++
			k = key + strconv.Itoa(n)
			_, ok = CM.JobGroups[k]
			continue
		}
		key = k
		break
	}

	CM.JobGroups[key] = group
	return group
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

//func (CM *copierMaschine) groupExists(jobName string) bool {}

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

// filecopy acts as a record of a single file's copy operation
type filecopy struct {
	relpath         string
	inSize, outSize int64
}

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
