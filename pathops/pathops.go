package pops

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func autoerr(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Panic(e)
		}
	}
}

func OpenReadall(fpath string) *os.File {
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
func makeOpenFileF(fpath string) *os.File {
	e := os.MkdirAll(filepath.Dir(fpath), os.ModeDir)
	autoerr(e)
	file, e := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if os.IsExist(e) {
		file, e = os.Open(fpath)
		if e != nil {
			panic(fmt.Errorf("error: %w \ndatafile: %s exists but failed to open file", e, fpath))
		}
	}
	return file

}

func initConfig(dir string) {
	if dir == "" {
		dir = "~/dotget/"
	}

}

func myDir() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
