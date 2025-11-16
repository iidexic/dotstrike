package pops

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (J *CopyJob) builtinIgnores() {
	if J.BPrefs.IsOn(bNoRepo) {
		J.ignore.Patterns = append(J.ignore.Patterns, subptn{ptn: `.git`, matchDir: true})
		J.ignore.Patterns = append(J.ignore.Patterns, subptn{ptn: `.\.gitignore`, matchFile: true})
	}
	if J.BPrefs.IsOn(bNoHidden) {
		J.ignore.AddSubpattern(".*", true, true)
	}
}

// Run CopyJob using fs package walkdir
func (J *CopyJob) RunFS() error {
	if J.jobRan {
		return fmt.Errorf("CopyJob Already Ran")
	}

	J.builtinIgnores()
	if J.BPrefs[bRootSubdir] && J.parentPathOut == "" { // parentpath check to be safe
		J.parentPathOut = J.PathOut
		J.PathOut = Joinpath(J.PathOut, filepath.Base(J.PathIn))
	}
	J.jobRan = true
	df := os.DirFS(J.PathIn)
	// ── Walk ────────────────────────────────────────────────────────────
	e := fs.WalkDir(df, ".", J.WalkFS)
	if e != nil {
		J.OpErrors = append(J.OpErrors, fs.PathError{Path: J.PathIn, Err: e, Op: ""})
		if e == ErrNilDirEntry {
			return fmt.Errorf("(nil DirEntry) The Input path is not a direcory or does not exist")
		}
		return e
	}
	return nil
}

// ╭─────────────────────────────────────────────────────────╮
// │                      WALK FUNCTION                      │
// ╰─────────────────────────────────────────────────────────╯
func (J *CopyJob) WalkFS(p string, d DirEntry, e error) error {
	var inpath, outpath string
	if filepath.IsAbs(p) {
		return fmt.Errorf("Abspath FS paths not supported. Use Absolute FS & Local Root")
	}
	if d == nil {
		return ErrNilDirEntry
	}
	inpath = Joinpath(J.PathIn, p)
	outpath = Joinpath(J.PathOut, p)
	rec, err := J.record.newRecord(outpath, p, d.IsDir())
	if be := J.checkAndLogError(outpath, "newRecord (os.Stat)", err); be && J.abort {
		return err
	}
	if d.IsDir() { // ──────────────────────────────────────────────────────────
		e := J.walkPathDir(inpath, outpath) //walkPathDir logs error
		if e != nil && J.abort {
			return e
		}
		return nil
	} // ───────────────────────────────────────────────────────────────────────
	info, e := d.Info()
	J.checkAndLogError(inpath, fmt.Sprintf("%s DirPath.Info()", d.Name()), e)
	rec.setOrigin(info)
	//rec.setOriginalSize(info.Size())

	if J.ignore.isIgnored(p, false) {
		return nil
	}
	if J.BPrefs[bNoFiles] { // for dry runs
		return nil
	}
	//if looksLikeRawCopy(outpath, info) {}

	// ── 1. open in file ──
	inF, e := OpenExistingFile(inpath)
	if e != nil {
		J.logError(p, "OpenExisting_In", e)
		return nil //skip file if opening errors
	}
	defer inF.Close()
	// ── 2. make/open out file ──

	outF, e := ForceMakeFile(outpath)
	if J.checkAndLogError(outpath, "ForceMakeFile", e) {
		return e
	}
	defer outF.Close()

	// ── 3. perform copy ──
	wb, e := io.CopyBuffer(outF, inF, nil)
	//TODO: Pull an outF.Stat(), add modtime to rec
	if J.checkAndLogError(outpath, "CopyError_Out", e) {

	}
	// get modtime and send 2 rec
	iout, e := outF.Stat()
	if J.checkAndLogError(outpath, "Stat() error", e) {
	}
	rec.setNew(iout, wb)
	//rec.setNewSize(wb)
	if rec.sizeOrigin != rec.sizeFinal && !J.configCheck(bNoFiles) {
		J.logError(outpath, "CopySizeCheck", e)
	}
	return nil
}

func (J *CopyJob) walkPathDir(inPath, outPath string) error {
	// check ignore + prevent recursion (if PathOut is a subdir of PathIn)
	if J.ignore.isIgnored(inPath, true) || strings.HasPrefix(inPath, J.PathOut) {
		J.logDir(outPath, false)
		return fs.SkipDir
	}

	if J.BPrefs[bAllDirs] {
		e := os.MkdirAll(outPath, 0755) // MkdirAll(outPath,0) - no permission bits == no permissions
		if e != nil {
			J.logError(outPath, "walkPathDir-Mkdir", e)
			return e
		}
	}

	J.logDir(outPath, true)

	return nil
}
