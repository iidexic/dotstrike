package pops

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

// enum type for file op outcomes
type failureType int

func (f failureType) Detail() string {
	var rstr string
	switch f {
	case None:
		rstr = fmt.Sprintf("Path Operation successful")
	case BadPattern:
		rstr = "bad pattern provided"
	case DirNotExist:
		rstr = "directory does not exist"
	case FileNotExist:
		rstr = fmt.Sprintf("file does not exist")
	case FileExist:
		rstr = fmt.Sprintf("file already exists")
	case FailedOpen:
		rstr = fmt.Sprintf("file seems to exist, but failed to open")
	case Error:
		rstr = fmt.Sprintf("General Error")
	}
	return rstr

}

// Possible outcomes for attempting filesystem operations
// Different operations have the potential to trigger different subsets of these outcomes
const (
	None failureType = iota
	BadPattern
	DirNotExist
	FileNotExist
	FileExist
	PermissionDenied
	FailedOpen
	Error // if error is returned, an error will also be returend in PathActionResult.Err
)

type ReadResult struct {
	Contents []byte
	readpath string
	f        os.File
	Fail     failureType //just replace this with an error type and send back the error
	Err      error
}

func (rr ReadResult) OpPath() string { return rr.readpath }
func (rr ReadResult) Failed() bool   { return rr.Fail != None } //in use

/* func (rr *ReadResult) Explain() string {
	return rr.Fail.Detail(rr.readpath, rr.Err)
} */

//var errDie []error = []error{filepath.ErrBadPattern, path.ErrBadPattern, os.ErrExist, os.ErrNotExist,os.ErrPermission}

func ce(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Panic(e)
		}
	}
}

var Open = os.Open

// OpenExistingFile attempts to open an existing file
// on success: returns open *os.File, nil. on fail: returns nil, error
func OpenExistingFile(fpath string) (*os.File, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return file, err
	}
	return file, nil
}
func Result() {

}

// makeabs returns absolute path of inpath
// inpath may or may not be relative from home dir/cwd
func makeabs(inpath string) string {
	if !path.IsAbs(inpath) {
		var e error
		inpath, e = filepath.Abs(inpath)
		if e != nil {
			log.Panic(e)
		}
	}
	return inpath
}

// MakeOpenFileF will open the given fpath as a file. It will make the file if it does not exist,
//
//	and it will make any missing directories necessary.
func MakeOpenFileF(fpath string) *os.File {
	// Check 3 locationsthat
	e := os.MkdirAll(filepath.Dir(fpath), os.ModeDir)
	ce(e)
	file, e := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if os.IsExist(e) {
		file, e = os.Open(fpath)
		if e != nil {
			log.Panic(fmt.Errorf("error: %w \ndatafile: %s exists but failed to open file", e, fpath))
		}
	}
	return file
}

// ReadFile will read contents of file into a ReadResult object and return a ptr
// result contains file and/or operation outcome/error if e!=nil
func ReadFile(pathElements ...string) *ReadResult {
	fpath := path.Join(pathElements...)
	result := &ReadResult{Fail: None, readpath: fpath}
	file, e := os.ReadFile(fpath)
	if e != nil && os.IsNotExist(e) {
		result.Fail = FileNotExist
		result.Err = e
	} else if e != nil {
		result.Fail = FailedOpen
		result.Err = e
	}
	if e != nil {
		result.Err = e

		if os.IsNotExist(e) {
			result.Fail = FileNotExist
		} else if os.IsPermission(e) {
			result.Fail = PermissionDenied
		} else {
			result.Fail = Error
		}
	}
	result.Contents = file //? cause panic on failure to read?
	return result
}

//TODO: migrate to ReadFileOrErr method; remove error wrapper

// ReadFileOrErr will read contents of file into a ReadResult and return a ptr to it.
// adds error to result.Err if not nil. Does not populate result.Fail
func ReadFileOrErr(pathElements ...string) *ReadResult {
	fpath := path.Join(pathElements...)
	result := &ReadResult{Fail: None, readpath: fpath}
	file, e := os.ReadFile(fpath)
	if e != nil {
		result.Err = e
	}
	result.Contents = file
	return result
}

func CalledFrom() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		//what
	}
	return dir
}

// TODO: symlink testing (next 2 functions related)
func IsBasicPath(p string) bool {
	return filepath.IsAbs(p) || filepath.IsLocal(p) //No symlink check/condition right now.

}
func IsSymlink(p string) bool {
	symP, e := filepath.EvalSymlinks(p)
	ce(e)
	//TODO: symlink testing
	return symP == p
}

// CheckPath for debug. TODO:remove when done
func CheckPath(p string) string {
	abs, e := filepath.Abs(p)
	ce(e)
	isabs := filepath.IsAbs(p)
	// abs: `//` or `\\`
	// loc:  (letters-only)||(starts with .)||(non-legal shit like '&', '^')
	// NEITHER: - `` (backticks - empty str)||(starts with single backslash or forward slash)||:(colon)

	/* NOTE:
	1. what needs to be handled?
		- forward slash/back slash
	*/
	isloc := filepath.IsLocal(p)
	var ptypestr string
	if isabs {
		ptypestr = ptypestr + "[absolute]"
	}
	//why not check for both
	if isloc {
		ptypestr = ptypestr + "[local]"
	}
	if !isabs && !isloc {
		ptypestr = ptypestr + "[UNKNOWN]"
	}
	path.Clean(p)

	locpls := filepath.Clean(p)
	base := filepath.Base(p)
	dir := filepath.Dir(p)
	return fmt.Sprintf(`
---| Check Path: '%s'
---------
abs:%t | local:%t

make abs path: 
	%s
clean/shorten path:
	%s
base = %s, dir = %s
`, p, isabs, isloc, abs, locpls, base, dir)

}
