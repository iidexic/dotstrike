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

type failureType int

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
	opfail(failureType, error)
	explain() string
	OpPath() string
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

func (rr *ReadResult) opfail(t failureType, e error) {
	rr.Fail = t
	rr.Err = e
}
func (rr ReadResult) OpPath() string { return rr.readpath }
func (rr *ReadResult) explain() string {
	var rstr string
	switch rr.Fail {
	case None:
		rstr = fmt.Sprintf("No failure type. Data read:\n %s", string(rr.Contents))
	case BadPattern:
		rstr = "bad pattern provided"
	case DirNotExist:
		rstr = "directory does not exist"
	case FileNotExist:
		rstr = fmt.Sprintf("file %s does not exist", rr.readpath)
	case FileExist:
		rstr = fmt.Sprintf("file %s already exists", rr.readpath)
	case FailedOpen:
		rstr = fmt.Sprintf("file %s seems to exist, but failed to open", rr.readpath)
	case Error:
		rstr = fmt.Sprintf("General Error: %e\nPath: %s", rr.Err, rr.readpath)
	}
	return rstr
}

//var errDie []error = []error{filepath.ErrBadPattern, path.ErrBadPattern, os.ErrExist, os.ErrNotExist,os.ErrPermission}

func ce(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Panic(e)
		}
	}
}

// OpenFile(fpath) opens a file; whether or not it or its parent directories exist
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
//  and it will make any missing directories necessary.
/* Basically, it will tear its way to whatever you give it, even if it doesn't exist */
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
func ReadFile(fpath string) *ReadResult {
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
