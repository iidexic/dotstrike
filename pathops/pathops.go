package pops

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

type fileReadResult struct {
	Contents any
	Fail     bool
}

func ce(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Panic(e)
		}
	}
}

// OpenFile(fpath) opens an existing file
func OpenFile(fpath string) *os.File {
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

// TY says ur welcome
func TY() {
	print("yw :)")
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
func ReadF(fpath string) *fileReadResult {
	file, e := os.ReadFile(fpath)
	if e != nil && os.IsNotExist(e) {
		return &fileReadResult{Fail: true}
	} else if e != nil {
		panic(fmt.Errorf("error: %w \ndatafile: %s exists but failed to open file", e, fpath))
	}

	return &fileReadResult{Contents: file}

}

func myDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
