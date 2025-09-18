package pops

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type FS = fs.FS
type File = fs.File
type FileInfo = fs.FileInfo
type FileMode = fs.FileMode
type uniFS map[string]uniFile

type uniFile struct {
	mode FileMode
	info FileInfo
}

// ── Test FS Implementation ──────────────────────────────────────────
func (F uniFile) Stat() (_ fs.FileInfo, _ error) {
	panic("not implemented") // TODO: Implement
}
func (F uniFile) Read(_ []byte) (_ int, _ error) {
	panic("not implemented") // TODO: Implement
}
func (F uniFile) Close() (_ error) {
	panic("not implemented") // TODO: Implement
}

func (U uniFS) Open(name string) (File, error) {
	f, ok := U[name]
	if ok {
		return f, nil
	}
	return nil, fs.ErrNotExist
}

func makeTestFS() *uniFS {
	return nil
}

// ──────────────────────────────────────────────────────────────────────

func TestPathExplore(t *testing.T) {
	pfile := "D:/coding/exampleFiles/INPUT/file_format/bad-gif.gif"
	split := filepath.SplitList(pfile)
	clean := filepath.Clean(pfile)
	t.Logf("clean:%s\nsplit:%v", clean, split)
	localize, err := filepath.Localize(filepath.Clean(pfile))
	t.Logf("Localized = %s (err %v)", localize, err)
}

func TestPathTear(t *testing.T) {
	usepath := "D:/coding/exampleFiles/INPUT"
	t.Logf("Split path: %s", usepath)
	listp := splitPath(usepath)
	t.Log("List Path:")
	for i, bp := range listp {
		t.Logf("[%d] %s", i, bp)
	}
	t.Logf("recombined: %s", filepath.Join(listp...))
}

func testing_job() *CopyJob {
	in := CleanPath(`d:\coding\exampleFiles\INPUT\`)
	out := CleanPath(`d:\coding\exampleFiles\OUTPUT`)
	job := CopyJob{
		PathIn: in, PathOut: out,
		newDirs: make(map[string]bool),
		BPrefs:  make(boolConfig),
		SPrefs:  make(stringConfig),
	}
	job.BPrefs[bNoRepo] = true
	job.BPrefs[bUseGlobal] = true
	job.BPrefs[bRootSubdir] = true
	job.BPrefs[bNoFiles] = false
	job.BPrefs[bNoHidden] = false

	return &job

}

func TestRunFSdry(t *testing.T) {
	job := testing_job()
	job.BPrefs[bNoFiles] = true
	e := job.RunFS()
	if e != nil {
		t.Errorf("JobRun Error: %v", e)
	}
	t.Log(job.Detail())
	t.Log(job.DetailRunFiles())
	if len(job.OpErrors) > 0 {
		t.Log("OPERATION ERRORS:")
		for _, e := range job.OpErrors {
			t.Logf("%v", e)
		}
	}
}

func TestRunFSDirs(t *testing.T) {
	job := testing_job()
	dl, e := newDirLog(job)
	if e != nil && dl != nil {

	}
	job.BPrefs[bNoFiles] = true
	e = job.RunFS()
	if e != nil {
		t.Errorf("JobRun Error: %v", e)
	}

	for _, e := range job.OpErrors {
		t.Logf("%v", e)
	}
}

func TestDirFS(t *testing.T) {
	pathdir := CleanPath(`d:/coding/exampleFiles/INPUT`)

	t.Logf("DirFS at '%s'", pathdir)
	d := os.DirFS(pathdir)
	opath := "."
	f, e := d.Open(opath)
	if e != nil {
		t.Errorf("Opening '%s' Does not work.", opath)
	}
	t.Logf("Opened File:%+v", f)
	defer f.Close()
	fstat, e := f.Stat()
	if e != nil {
		t.Errorf("Stat @ '%s' failed: %v", opath, e)
	}
	t.Logf("fstat: %+v", fstat)

}

func TestDirFSWalk(t *testing.T) {
	pathdir := CleanPath(`d:/coding/exampleFiles/INPUT`)

	d := os.DirFS(pathdir)
	opath := "."
	e := fs.WalkDir(d, opath, func(p string, d DirEntry, e error) error {
		if e != nil {
			t.Logf("e @ start of walk: %s", e.Error())
		}
		info, e := d.Info()
		if e != nil {
			t.Logf("(error on Info(): %v)", e)
		}
		t.Logf("Walk: p='%s', isAbs:%t, isDir:%t, Type:%+v, Size:%v",
			p, filepath.IsAbs(p), d.IsDir(), d.Type(),
			float64(info.Size()/100000.0))
		return nil
	})
	if e != nil {
		t.Errorf("Walk Error: %v", e)
	}

}
