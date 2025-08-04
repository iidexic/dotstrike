package pops

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func getcwd(t *testing.T) string {
	cwd, e := os.Getwd()
	if e != nil {
		t.Logf("[cwd: %e]", e)
	}
	return cwd
}

func TestRead(t *testing.T) {
	lpath := "../_xtra/dotstrike.toml"
	abspath, e := filepath.Abs(lpath)
	readout := ReadFile(lpath)
	if e != nil {
		t.Logf("Abspath Error:%e", e)
	}
	t.Logf("Path: %s", readout.OpPath())
	t.Logf("Abs Path: %s", abspath)
	t.Log(string(readout.Contents))
	t.Fail()
}

func TestMakeAbs(t *testing.T) {
}

func TestPathIdentification(t *testing.T) {
	trials := []string{
		"``", `\`, "''", "_", "=",
	}
	_ = trials
}

func TestScratch(t *testing.T) {
	l := []string{"a", "two", "five"}
	t.Logf("l: %s", l)
	t.Logf("l as v: %v", l)
	lc := slices.Clone(l)
	t.Logf("l-clone: %s", l)
	t.Logf("original match clone? -> %t", slices.Equal(l, lc))

	lc = lc[:]
	t.Logf("original match clone[:]? -> %t", slices.Equal(l, lc))
	l = l[:0]
	clear(lc)
	t.Logf("l[:0]: %s", l)
	t.Logf("clear(lclone): %s (length==%d)", lc, len(lc))
	t.Logf("[:0] match clear()? -> %t", slices.Equal(l, lc))
	t.Log("[:0] wipes length, clear does not")
	t.Fail()
}

type trojhrse struct {
	pl      []string
	plsize  []int64
	ignores []string
	t       *testing.T
	n       int
}

func (th *trojhrse) walkF(p string, d fs.DirEntry, e error) error {
	th.n++
	if fname := d.Name(); strings.Index(fname, ".") == 0 || strings.Index(fname, "_") == 0 {
		if d.IsDir() {
			th.t.Logf("skipdir: %s. Details: Type()=%v\n(isregular?->%t)", fname, d.Type(), d.Type().IsRegular())

		}
		//Create new paths needed for copy here probably
		return nil
	}

	if gi := strings.Index(filepath.Dir(p), ".git"); gi >= 0 {
		return nil
	}
	if e != nil {
		th.t.Logf("path %s: error %V", p, e)
		return e
	}
	fi, eloc := d.Info()
	if eloc != nil {
		th.t.Logf("error %s.Info(): %e", d.Name(), eloc)
		return eloc
	}
	th.t.Log()

	th.pl = append(th.pl, p)
	th.plsize = append(th.plsize, fi.Size())
	return nil
}

var th = trojhrse{}

func TestWalkwd(t *testing.T) {
	cwd := filepath.Dir(getcwd(t))
	t.Log("WD =", cwd)
	th.t = t
	ew := filepath.WalkDir(cwd, th.walkF)

	if ew != nil {
		t.Errorf("Walk error\n[[%e]]", ew)
	}
	t.Logf("Details of Walk:\nQty Paths:%d, qty size:%d", len(th.pl), len(th.plsize))
	if len(th.pl) > 0 {
		for i, v := range th.pl {
			t.Logf("%d) - %s", i, v)
		}
	}
}

func TestCleanPath(t *testing.T) {
	paths := []string{"../_xtra/dotstrike.toml", ".\\path partee\\\\", "./dotcheck/./doit.go",
		"..\\..\\double_dot_check\\..\\also-slashes\\", "junktext:P*#!%H!PO( TVHL)}\n", "~/tilde?"}
	for i, p := range paths {
		cleaned := CleanPath(p)
		t.Logf("Path [%d]:\n (%s)->(%s)", i, p, cleaned)
	}
}
