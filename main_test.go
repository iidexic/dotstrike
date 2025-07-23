package main

import (
	"os"
	"path/filepath"
	"testing"
)

func getcwd(t *testing.T) string {
	cwd, e := os.Getwd()
	if e != nil {
		t.Logf("[cwd: %e]", e)
	}
	return cwd
}
func retAsIs(fpath string) error {
	f, e := os.OpenFile(fpath, os.O_RDONLY, 0)
	f.Close()
	return e
}
func retCheck(fpath string) error {
	f, e := os.OpenFile(fpath, os.O_RDONLY, 0)
	f.Close()
	if e != nil {
		return e
	}
	return nil
}
func TestErrorReturn(t *testing.T) {
	echek := retCheck("./notarealfile.fileextension")
	t.Logf("Checked done. Got: %v", echek)
	enochek := retAsIs("./notarealfile.fileextension")
	t.Logf("Unchecked done. Got: %v", enochek)
	if echek != enochek {
		t.Logf("is echek==enochek? %t", echek == enochek)
		t.Logf("(same error type but different instances of error)")
	}
	t.Logf("echek = no file? %t\nenochek = no file? %t", os.IsNotExist(echek), os.IsNotExist(enochek))
}
func TestStringIndexing(t *testing.T) {
	str := "~\\what/$&$#^"
	for i := range len(str) {
		t.Logf("%d) %d %s", i, str[i], string(str[i]))
	}
}

func TestPathChanges(t *testing.T) {
	cwd := getcwd(t)
	lop := func(oname, res string) {
		t.Logf("|%s()-> %s", oname, res)
	}
	pdir := filepath.Dir(cwd)
	gpdir := filepath.Dir(pdir)
	t.Logf("cwd = %s", cwd)
	lop("Base", filepath.Base(cwd))
	lop("Dir", pdir)
	lop("Dir x2", gpdir)
	lop("Clean", filepath.Clean(cwd))
	lop("Ext", filepath.Ext(cwd))
	lop("VolumeName", filepath.VolumeName(cwd))
	lop("FromSlash", filepath.FromSlash(cwd))
	dirr, fll := filepath.Split(cwd)
	t.Logf("|Split()-> %s, %s", dirr, fll)
	gfiles, e := filepath.Glob("*.*")
	if e != nil {
		t.Logf("[glob: %e]", e)
	}
	t.Logf("GLOB:\n------")
	for i, fn := range gfiles {
		t.Logf("(%d) %s", i, fn)
	}

}
func TestPathSplit(t *testing.T) {
	cwd := getcwd(t)
	cwdloc, e := filepath.Localize(filepath.FromSlash(filepath.Clean(cwd)))
	if e != nil {
		t.Logf("localize err: %e", e)
	}
	seg := filepath.SplitList(cwdloc)
	for i, s := range seg {
		t.Logf("%d. %s", i, s)
	}
}
