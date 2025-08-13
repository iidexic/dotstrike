package pops

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var dataTx copyProcessing

func (J *CopyJob) collectDirs() chan<- pathcopyStatus {
	ch := make(chan pathcopyStatus, 64)
	go func() {
		for d := range ch {
			J.logDir(d.path, d.copied)
		}
	}()
	return ch
}
func (J *CopyJob) RunC( /* params */ ) error {
	// condition paths
	var e error

	if J.JobSettings.makeRootSubdir && J.parentPathOut == "" { // parentpath check to be safe
		J.parentPathOut = J.PathOut
		J.PathOut = Joinpath(J.PathOut, filepath.Base(J.PathIn))
	}
	// ── ensure chan struct loaded ───────────────────────────────────────
	dataTx.cdirs = J.collectDirs()
	dataTx.cfp = make(chan fs.PathError, 64)

	// ── Walk ────────────────────────────────────────────────────────────
	e = filepath.WalkDir(J.PathIn, J.Walk)

	if e != nil {
		J.OpErrors = append(J.OpErrors, fs.PathError{Path: J.PathIn, Err: e, Op: ""})
		return e
	}

	//WARNING: WITH THIS STRUCTURE, A WALKDIR ERROR WILL PREVENT MAKING ADDITIONAL DIRS
	if J.JobSettings.copyAllDirectories {
		for dir := range J.newDirs {
			e := os.MkdirAll(Joinpath(J.PathOut, dir), 0)
			J.checkAndLogError(dir, "MakeDirectory", e)
		}
	}
	return nil
}

func processDir(path string, written bool) {
	defer wg.Done()

}
func fileCCopy(inpath, outpath string) {

}

func (J *CopyJob) dirlogC(relpath string, written bool) {

}
func (J *CopyJob) logErrors(perrs ...fs.PathError) {

}

type copyProcessing struct {
	cdirs chan<- pathcopyStatus
	cfp   chan<- fs.PathError
}
type pathcopyStatus struct {
	path   string
	copied bool
}

func (J *CopyJob) WalkC(p string, d DirEntry, e error) error {
	// make relative path first; used for dirs & files
	prr := J.stripRoot(p) //	!INFO: panics on error; error is unexpected

	// DIRECTORIES:
	if d.IsDir() {
		var re error
		writedir := true
		// check ignore + prevent recursion (if PathOut is a subdir of PathIn)
		if J.ignore.isIgnored(p, true) || strings.HasPrefix(p, J.PathOut) { // WARN: need to check for files too
			writedir = false
			re = fs.SkipDir
		}
		wg.Add(1)
		go func() { //pipe dir record over
			defer wg.Done()
			dataTx.cdirs <- pathcopyStatus{prr, writedir}
		}()
		return re
	}

	// FILES:
	// ── 0.1 make filepath out ─────────────────────────────────
	pto := Joinpath(J.PathOut, prr)
	if !filepath.IsAbs(pto) {
		pto = absNoE(pto) //	!INFO: panics on error; error is unexpected
	}

	// 0.2 Get infile Info()
	inDE, e := d.Info()
	J.checkAndLogError(p, "GetFileInfo_In", e)

	// 0.3 Return before copy if config requires.
	if J.JobSettings.noFiles {
		J.addFile(prr, inDE.Size(), 0) // for dry runs
		return nil
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
	J.addFile(prr, inDE.Size(), wb)
	return nil
}
