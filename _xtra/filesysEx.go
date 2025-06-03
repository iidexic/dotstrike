//package pops

/*
import (
	//	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func dig() {
	contents, e := os.ReadDir(".")
	ce(e)
	for i, c := range contents {
		fmt.Print(i, ":", c.Name(), " - ", filepath.Ext(c.Name()))
	}
}

// OpenFile opens a file while working out potential errors as best I know how
func openFile(fpath string) *os.File {
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

func storageFile(pathf string) {
	file := openFile(pathf)
	// need to handle err on file.close. jk dont need another function
	// need to move/remove if working with file in a diff function?
	defer ce(file.Close())
	// bingo now do what you gotta do with this file
	finf, e := file.Stat()
	ce(e)
	content := make([]byte, finf.Size())
	rnum, e := file.Read(content)

	fmt.Print(rnum)
	ce(e)
}

func readF(openfile os.File) {
} */

/* now have all the tools to build the crud dotget
Filesystem Manipulation Notes
=============================
__1. Read Files__
	`os.ReadFile(str)` to just read a whole file
for more nuanced reads:
	`file := os.Open(str)` to get an `os.File`
	`file.Read([]byte) int, err` will read bytes into the byteslice, up to `len([]byte)`, returns # read
	`file.Seek(int, io."Position")` will move the read cursor by +(int) positions relative to io."Position"
Positions from io:
	-> io.SeekStart = beginning of file
	-> io.SeekCurrent = cursor location
	-> io.SeekEnd = end of file
Other method of file.Read()
	io.ReadAtLeast(os.File, []byte, int) int, err :: reads at least [int], at most [[]byte] size.
__2. Write Files__
	`os.WriteFile(filestr, (data? as []byte), integer?)
for more nuanced writes:
	file,err := os.Create(filestr)
	check(err)
	defer file.Close()
	file.Write([]byte) int, err : []byte = data to write, return int = bytes written
	file.WriteString("text to write") int err : same but direct string write

	file.Sync() <- issue a Sync to flush writes to stable storage (write to disk, from memory)

For both reads and writes, the bufio package has buffered reader/writer

__3. File Paths__
using 'path/filepath' package:
Join() used  to construct paths. Always do this instead of regular string concatenation
	`p:=filepath.Join(dir1str,dir2str,filenamestr)
	filepath.Dir(p) -> split path to just directory
	filepath.Base(p) -> split path to just filename
	filepath.Split(p) -> split path and return both directory, filename
		also filepath.IsAbs() -> bool if absolute path
File extensions:
	extension:=filepath.Ext(filename)
	name_only:=strings.TrimSuffix(filename,ext) -> removes file extension giving only name
relative path
	rel, e := filepath.Rel(basepathstr, targetpathstr) -> i.e. 'c:/docs','c:/docs/projects/p1.txt'
__4. Directories_
	os.Mkdir("subdir", 0755)-> runs mkdir RELATIVE TO CWD
	defer os.RemoveAll("subdir") -> delete dir+contents, use defer when making temp files
Q: what is 0755 for? 0755==mkdir? 0644==writefile???
make file:
	check(os.WriteFile(name, []byte (data), 0644))->for blankfile, data = []byte("")
make dir multi:
	os.MkdirAll("subdir/parent/child", 0755) -> makes all dirs needed to get to have child dir (FROM CWD)
dir move/read:
	os.Chdir("subdir/parent/child") navigates CWD to given path
	os.ReadDir(".") -> reads dir contents and returns. (rel 2 CWD? prob)
Walk:
	filepath.WalkDir("subdir", visit) -> recursively visit a directory/subdirs
		here visit = callback function we can write/provide, handles each file or dir visited
__4. Temp Files__
f, err := os.CreateTemp("","sample") -> create temp file (location,""=default for os, name??)
	defer os.Remove(f.Name()) clean up directly (rel 2 CWD or abspath)
td, e := os.MkdirTemp("","sampledir") -> create temp directory
	defer os.RemoveAll(td)
now use temp dir to build temp filenames:
	fname := filepath.Join(td, filename)
*/
