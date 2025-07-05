package pops

import (
	"io/fs"
	"path/filepath"
)

type DirEntry = fs.DirEntry

// CopierMaschine builds and executes CopyJobs
type CopierMaschine struct {
	JobQueue []CopyJob
}

// CopyJob prepares and executes the copy of all contents of PathIn to PathOut
type CopyJob struct {
	PathIn, PathOut string
	fstack          []filecopy
	fsCopied        []float32
	ignore          IgnoreSet
}

// filecopy will run a single file copy task and acts as a record of that task
// TODO:re-tool; Originally did not plan to perform all copying in a single WalkDir
type filecopy struct {
	name            string
	id              int //?
	inPath, outPath string
	inSize, outSize int
}

func (J *CopyJob) Load(absPathIn, absPathOut string) {
	J.fstack = append(J.fstack, filecopy{
		inPath: absPathIn, outPath: absPathOut,
	})
}

func (J *CopyJob) newFcopy(in, out string) *filecopy {
	return &filecopy{inPath: in, outPath: out}
}

func (J *CopyJob) Walk(p string, d DirEntry, e error) error {
	if d.IsDir() {
		if J.ignore.IsIgnoredDir(p) {
			return fs.SkipDir
		}
		return nil
	}

	return nil
}

func (J *CopyJob) RunCopyWalk() error {
	e := filepath.WalkDir(J.PathIn, J.Walk)
	if e != nil {
		return e
	}
	return nil
}

// ── Ignore Functionality ────────────────────────────────────────────
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

func (ip iptn) matches(s string) bool {
	return false
}

func (I *IgnoreSet) IsIgnoredDir(dirpath string) bool {
	for _, pat := range I.ignoreDir {
		if pat.matches(dirpath) {
			return true
		}
	}
	return false
}
