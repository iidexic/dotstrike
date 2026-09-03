package virtfs

import (
	"io"
	"io/fs"
	"sync"
	"time"
)

type (
	FS          = fs.FS
	File        = fs.File
	FileInfo    = fs.FileInfo
	ReadDirFS   = fs.ReadDirFS
	ReadDirFile = fs.ReadDirFile
	dirEntry    = fs.DirEntry
)
type uniFS struct {
	files map[string]uniFile
}

type rwtree interface {
	write(string, []byte) error
}

const (
	modeDir    = fs.ModeDir
	modeDevice = fs.ModeDevice
)

//TODO:  In future, add cursor/position to files for both read/write

type uniTree map[string]*uniFile

type ubasefile struct {
	name string
	mu   *sync.Mutex
	mode fs.FileMode
	info uniInfo
	mod  time.Time
}

type uniFile struct {
	ubasefile
	data []byte
}

type uniDir struct {
	uniFile
	tree uniTree
}

func (uniDir) ReadDir(n int) (_ []fs.DirEntry, _ error) {
	panic("not implemented") // TODO: Implement
}

type uniInfo struct {
	file *uniFile
}

func makemode()

// ── Test FS Implementation ──────────────────────────────────────────
func (F uniFile) Stat() (_ fs.FileInfo, _ error) {
	return nil, nil
}

// appends uniFile contents to b.
// errors if
func (F uniFile) Read(b []byte) (_ int, _ error) {
	if F.mode == modeDir {
		return 0, io.EOF
	}
	if ldat := len(F.data); ldat > 0 {
		b = append(b, F.data...)
		return ldat, nil
	}
	return 0, nil
}

func (F uniFile) Close() (_ error) {
	F.mu.Unlock()
	panic("not implemented") // TODO: Implement
}

func (U uniFS) Open(name string) (File, error) {
	f, ok := U.files[name]
	if ok {
		return f, nil
	}
	return nil, fs.ErrNotExist
}

func (U *uniFS) OpenRW(name string) (File, error) {
	f, ok := U.files[name]
	if ok {
		return f, nil
	}
	return nil, fs.ErrNotExist
}

func (U *uniFS) Write(file string, data []byte) (_ error) {
	f, ok := U.files[file]
	if !ok {
		return fs.ErrNotExist
	}
	if f.mode.IsDir() {
		return fs.ErrInvalid
	}
	copy(f.data, data)
	return nil
}

func (I uniInfo) Name() (_ string)       { return I.file.name }
func (I uniInfo) Size() (_ int64)        { return int64(len(I.file.data)) }
func (I uniInfo) Mode() (_ fs.FileMode)  { return I.file.mode }
func (I uniInfo) ModTime() (_ time.Time) { return I.file.mod }
func (I uniInfo) IsDir() (_ bool)        { return I.file.mode.IsDir() }
func (I uniInfo) Sys() (_ any)           { return "virtual" }
