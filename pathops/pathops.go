package pops

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// enum for path format
type pathType int16

const (
	UnknownPath pathType = iota - 1
	InaccessiblePath
	//---
	LocalDirPath
	AbsDirPath
	LocalFilePath
	AbsFilePath
)

var HomePath *string = nil

type errCtr struct {
	etext string
	err   error
}

func (e *errCtr) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: [%s]", e.etext, e.err.Error())
	}
	return e.etext
}

var ErrGetHome = errCtr{etext: "Failed to retrieve user homedir"}
var ErrEmptyHome = fmt.Errorf("Home path is empty string")
var ErrNilInfo = fmt.Errorf("nil os.FileInfo")

var Open = os.Open
var BaseName = filepath.Base

// Joinpath aliases filepath.Join (no longer necessary)
var Joinpath = filepath.Join

func ce(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			panic(e)
		}
	}
}

//TODO:(low-recl) Replace ALL os.IsExist/os.IsNotExist with errors.Is()
//TODO:(med-recl) Clean up Home functions - here and where used
//TODO:(med-feat) Replace current config path system with more robust system with fallbacks

// enum type for file op outcomes
type failureType int16 //TODO:(med-recl) replace system with errors; currently basically converting errors to enum

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
const (
	tilde     = byte('~')
	backslash = byte('\\') //unused
	amp       = byte('&')  //unused
	atsign    = byte('@')  //unused
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
return rr.Fail.Detail(rr.readpath, rr.Err) } */

func (f failureType) Detail() string {
	var rstr string
	switch f {
	//TODO: check later to see if I actually ever use the fmt.Sprintf
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

// OpenExistingFile attempts to open an existing file
// on success: returns open *os.File, nil. on fail: returns nil, error
// NOTE: Read-Only
func OpenExistingFile(ospath string) (*os.File, error) {
	file, err := os.Open(ospath)
	if err != nil {
		return file, err
	}
	return file, nil
}

// OpenFileRW wraps os.OpenFile(ospath, os.O_RDWR,0)
func OpenFileRW(ospath string) (*os.File, error) {
	return os.OpenFile(ospath, os.O_RDWR, 0)
	//Q: need check error before?
}

// ── HOME PATH FUNCTIONS ─────────────────────────────────────────────

// HomeJoin retrieves abs homedir path, adds suffix to the end, and returns.
// Directly returns error from os.UserHomeDir().
func HomeJoin(suffix string) (string, error) {
	if HaveHome() {
		return HomeJoinC(suffix), nil
	}
	home, e := os.UserHomeDir()
	if e != nil {
		return "", e
	}
	HomePath = &home
	return Joinpath(home, suffix), nil
}

// HomeJoinC uses HomePath var (populated on init) to prepend homedir to suffix
func HomeJoinC(suffix string) string { return Joinpath(*HomePath, suffix) }

// HomeDirtyJoin retrieves abs homedir path, adds suffix to the end, and returns.
// errors will panic
func HomeDirtyJoin(suffix string) string {
	home, e := os.UserHomeDir()
	ce(e)
	return filepath.Join(home, suffix)
}
func GetHomeDir() error {
	home, e := os.UserHomeDir()
	if e != nil {
		ErrGetHome.err = e
		return &ErrGetHome
	}
	if home != "" {
		HomePath = &home
		return nil
	}
	return ErrEmptyHome
}

func HaveHome() bool {
	if HomePath != nil && filepath.IsAbs(*HomePath) {
		return true
	}
	return false
}

// TildeDirty replaces a leading ~ with home path using HomeDirtyJoin
// errors will panic.
func TildeDirty(ospath string) string {
	// tilde code: 126
	if c1 := ospath[0]; c1 == tilde && *HomePath != "" {
		return HomeJoinC(ospath[0:])
	} else if c1 == tilde {
		return HomeDirtyJoin(ospath[0:])
	}
	return ospath
}

func TildeFix(ospath string) (string, error) {
	// tilde code: 126
	if c1 := ospath[0]; c1 == tilde && *HomePath != "" {
		return HomeJoinC(ospath[1:]), nil
	} else if c1 == tilde {
		return HomeJoin(ospath[1:])
	}
	return ospath, nil
}

// makeabs returns absolute path of inpath
// inpath may or may not be relative from home dir/cwd
func MakeAbs(inpath string) string {
	if !filepath.IsAbs(inpath) {
		var e error
		inpath, e = filepath.Abs(inpath)
		if e != nil {
			panic(e)
		}
	} else {
		inpath = filepath.Clean(inpath)
	}
	return inpath
}
func PathExists(path string) (bool, error) {

	path = CleanPath(path)
	s, e := os.Stat(path)
	if e != nil && errors.Is(e, os.ErrNotExist) {
		return false, e
	} else if e != nil && errors.Is(e, os.ErrExist) {

	}
	if s != nil {
		return true, nil
	}
	return true, ErrNilInfo
}

func PathTypeIfExists(path string) (pathType, error) {
	path = CleanPath(path)
	s, e := os.Stat(path)
	if e != nil && errors.Is(e, os.ErrNotExist) {
		return UnknownPath, e
	} else if e != nil {

	}
	isdir := s.IsDir()
	isabs := filepath.IsAbs(path)
	switch {
	case isdir && isabs:
		return AbsDirPath, nil
	case isdir:
		return LocalDirPath, nil
	case isabs:
		return AbsFilePath, nil
	default:
		return LocalFilePath, nil
	}

}

var Abs = filepath.Abs

// MakeOpenFileF will open the given fpath as a file. It will make the file if it does not exist,
// and it will make any missing directories necessary.
func MakeOpenFileF(fpath string) (*os.File, error) {
	// makedir->create+open|if exists->open
	e := os.MkdirAll(filepath.Dir(fpath), os.ModeDir)
	ce(e)
	file, e := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
	if os.IsExist(e) {
		file, e = os.OpenFile(fpath, os.O_RDWR, 0)
		return file, e
	}
	return file, e
}

func NoError(errors []error) bool {
	for _, e := range errors {
		if e != nil {
			return false
		}
	}
	return true
}
func CopyFileME(toPath, fromPath string) []error {
	var e error
	//TODO: deal with these errors in a better way
	ers := make([]error, 7)
	//e0- getabs(to)
	//e1-getabs(from)
	//e2-OpenExisting(from)
	//e3-stat(from)
	//e4-MakeOpen(to)
	//e5-Copy(to/from)
	//e6-writtenBytes vs size(from) custom error
	if !filepath.IsAbs(fromPath) {
		fromPath, e = filepath.Abs(toPath)
		ers[0] = e
	}
	if !filepath.IsAbs(toPath) {
		toPath, e = filepath.Abs(toPath)
		ers[1] = e
	}
	// 2. open fromPath
	fromFile, e := OpenExistingFile(fromPath)
	ers[2] = e
	defer fromFile.Close()
	fromStat, e := fromFile.Stat()
	ers[3] = e
	fromSize := fromStat.Size()

	toFile, e := MakeOpenFileF(toPath)
	ers[4] = e
	defer toFile.Close()

	writtenBytes, e := io.CopyBuffer(toFile, fromFile, nil)
	ers[5] = e
	if fromSize != writtenBytes {
		if writtenBytes > fromSize {
			ers[6] = fmt.Errorf("wrote more bytes (%d) than original filesize(%d)",
				writtenBytes, fromSize)
		} else {
			ers[6] = fmt.Errorf("Did not write enire file: size %d, written %d", fromSize, writtenBytes)
		}
	}
	return ers
}

// ReadFile will read contents of file into a ReadResult object and return a ptr
// result contains file and/or operation outcome/error if e!=nil
func ReadFile(pathElements ...string) *ReadResult {
	fpath := filepath.Join(pathElements...)
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
	fpath := filepath.Join(pathElements...)
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
		panic(err)
		//what
	}
	return dir
}

// CleanPath returns the shortest path name equivalent to path
// via lexical processing.
// The following rules are applied iteratively:
//
//  1. Replace multiple [Separator] elements with a single one.
//  2. Eliminate each . path name element (the current directory).
//  3. Eliminate each inner .. path name element (the parent directory)
//     along with the non-.. element that precedes it.
//  4. Eliminate .. elements that begin a rooted path:
//     that is, replace "/.." by "/" at the beginning of a path,
//     assuming Separator is '/'.
//
// The returned path ends in a slash only if it represents a root directory,
// such as "/" on Unix or `C:\` on Windows.
//
// Finally, any occurrences of slash are replaced by Separator.
//
// If the result of this process is an empty string, Clean
// returns the string ".".
var CleanPath = filepath.Clean

// dirContents returns a map describing contents of dirpath and subdirectories
// although not in very clear detail and it probably should not be used
// map keys are paths of existing filesystem objects
// value bool indicates whether or not that object is a directory
// it returns nil if - dir doesn't exist, path is not a dir, any os.Stat error, any WalkDir error
func DirContents(dirpath string) *map[string]bool {
	dirstat, e := os.Stat(dirpath)
	if e != nil || !dirstat.IsDir() {
		// no reason to check IsNotExists; we can only return nil
		return nil
	}
	contentsIsDir := make(map[string]bool)
	e = filepath.WalkDir(dirpath, func(p string, d DirEntry, e error) error {
		if d.IsDir() {
			contentsIsDir[p] = true
		}
		contentsIsDir[p] = false
		return nil
	})
	if e != nil {
		return nil
	}
	return &contentsIsDir
}

// MaybePasicPath just checks if p IsAbs or IsLocal
// So it probably should not be used
func MaybeBasicPath(p string) bool {
	return filepath.IsAbs(p) || filepath.IsLocal(p) //No symlink check/condition right now.

}
func IsSymlink(p string) bool {
	symP, e := filepath.EvalSymlinks(p)
	ce(e)
	//TODO: symlink testing
	return symP == p
}

// CheckPathDebug for debug. TODO:remove when done
func CheckPathDebug(p string) string {
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
	filepath.Clean(p)

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
