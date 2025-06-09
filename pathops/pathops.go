package pops

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"syscall"
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
	FailedOpen
	Error // if error is returned, an error will also be returend in PathActionResult.Err
)

type PathEvent interface {
	// opfail(failureType, error)
	Explain() string // returns result printable to user
	OpPath() string  // returns read path/path used
}
type MakeOpenResult struct {
	PathsMade []string
	Fail      failureType
	Err       error
}

type ReadResult struct {
	Contents []byte
	readpath string
	f        os.File
	Fail     failureType //just replace this with an error type and send back the error
	Err      error
}

// I don't actually know if these are needed
/* type OpResult interface{ MakeOpenResult | ReadResult }

func IsResErr[R OpResult](e error, result R) {}
func (rr *ReadResult) opfail(t failureType, e error) {
	rr.Fail = t
	rr.Err = e
}
*/

func (rr ReadResult) OpPath() string { return rr.readpath }

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

// OpenFile(fpath) opens a file; whether or not it or its parent directories exist
// TODO: OPENFILE AND MAKEOPENFILE ARE SUPPOSEDLY THE SAME BASED OFF OF DOC COMMENT.
//
//	Check differences; Remove OpenFileF or MakeOpenFileF, OR rewrite OpenFileF to only open if exists
func OpenFileF(fpath string) *os.File {
	file, err := os.Open(fpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOTDIR) {
			// not exist
			e := os.MkdirAll(filepath.Dir(fpath), os.ModeDir) // os.ModeDir right? check what is expected
			if e != nil {
				panic(e)
			}
			file, e = os.Create(fpath)
			if e != nil {
				panic(e)
			}
		} else {
			panic(err)
		}
	}
	return file
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
			panic(e)
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
			panic(fmt.Errorf("error: %w \ndatafile: %s exists but failed to open file", e, fpath))
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
