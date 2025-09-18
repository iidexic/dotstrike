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
	if J.BPrefs[bRootSubdir] && J.parentPathOut == "" { // parentpath check to be safe
		J.parentPathOut = J.PathOut
		J.PathOut = Joinpath(J.PathOut, filepath.Base(J.PathIn))
	}
	J.jobRan = true
	parent, base := filepath.Split(J.PathIn)
	if parent == "" || base == "" {
		return fmt.Errorf("Split did a bad job splitting")
	}
	df := os.DirFS(parent)
	// ── Walk ────────────────────────────────────────────────────────────
	e := fs.WalkDir(df, base, J.WalkFS)
	if e != nil {
		J.OpErrors = append(J.OpErrors, fs.PathError{Path: J.PathIn, Err: e, Op: ""})
		return e
	}
	// move to be inside the walk?
	if J.BPrefs[bAllDirs] {
		for relDir := range J.newDirs {
			e := os.MkdirAll(Joinpath(J.PathOut, relDir), 0)
			J.checkAndLogError(relDir, "MakeDirectory", e)
		}
	}
	return nil
}

// func MakeDirectories() {}

// TEST: BEFORE EVER USING THIS FUNC:
// func pathRelative(path string, parent string) string {
// 	if path == parent {
// 		return "."
// 	}
// 	if strings.Contains(path, parent) {
// 		if rn := rune(path[len(parent)]); rn == '/' || rn == '\\' {
// 			return strings.Replace(path, parent, ".", 1)
// 		}
// 	}
// 	path, parent = CleanPath(path), CleanPath(parent)
// 	if strings.Contains(path, parent) {
// 		return strings.Replace(path, parent, ".", 1)
// 	}
// 	return ""
// }

// ╭─────────────────────────────────────────────────────────╮
// │                      WALK FUNCTION                      │
// ╰─────────────────────────────────────────────────────────╯
func (J *CopyJob) WalkFS(p string, d DirEntry, e error) error {
	var inpath, outpath string
	if filepath.IsAbs(p) {
		return fmt.Errorf("Abspath FS paths not supported. Use Absolute FS & Local Root")
	}
	inpath = Joinpath(J.PathIn, p)
	outpath = Joinpath(J.PathOut, p)

	if d.IsDir() {
		e := J.walkPathDir(inpath, outpath)
		if e != nil && J.abort { //Logged in walkPathDir

			return e //need an abort check?
		}
	} else if J.ignore.isIgnored(p, false) {
		return nil
	}

	info, e := d.Info()
	J.checkAndLogError(p, fmt.Sprintf("%s DirPath.Info()", d.Name()), e)

	if J.BPrefs[bNoFiles] { // for dry runs
		J.addFile(p, info.Size(), 0)
		return nil
	}

	//if looksLikeRawCopy(outpath, info) {}

	// ── 1. open in file ──
	inF, e := OpenExistingFile(p)
	if e != nil {
		J.logError(p, "OpenExisting_In", e)
		return nil //skip file if opening errors
	}
	defer inF.Close()
	// ── 2. make/open out file ──

	outF, e := MakeOpenFileF(outpath)
	if J.checkAndLogError(outpath, "MakeOpen_Out", e) {
		outF.Close()
		return e
	}
	defer outF.Close()

	// ── 3. perform copy ──
	wb, e := io.CopyBuffer(outF, inF, nil)
	J.checkAndLogError(outpath, "CopyError_Out", e)

	// check size original matches new copied
	J.addFile(p, info.Size(), wb)
	return nil
}

func (J *CopyJob) walkPathDir(inPath, outPath string) error {
	// check ignore + prevent recursion (if PathOut is a subdir of PathIn)
	if J.ignore.isIgnored(inPath, true) || strings.HasPrefix(inPath, J.PathOut) {
		J.logDir(outPath, false)
		return fs.SkipDir
	}
	if J.BPrefs[bAllDirs] {
		e := os.MkdirAll(outPath, 0)
		if e != nil {
			J.logError(outPath, "walkPathDir-Mkdir", e)
			// TEST: check if returning the error terminates walkdir; if so, maybe add J.abort check?
			return e
		}
	}

	J.logDir(outPath, true)

	return nil
}
