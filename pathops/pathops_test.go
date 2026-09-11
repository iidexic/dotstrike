package pops

import (
	"io/fs"
	"os"
	"path/filepath"
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
func getHome(t *testing.T) string {
	home, e := os.UserHomeDir()
	if e != nil {
		t.Logf("[home: %e]", e)
	}
	HomePath = &home
	return home
}
func TestMakeAbs(t *testing.T) {
	trials := []string{
		"Bingo Bongo :)", "C:/Users/derek/.config/wezterm",
		".", "~", "~/.config/wezterm", "~/../loc/",
		"C:\\", "C\\", "c:/", "\\C\\",
		".\\bingo2", `D:\coding\exampleFiles`,
		"/", "\\", "..\\cmd", "..", "...", "....", "./", "./place/../things"}
	getHome(t)
	for i, p := range trials {
		abs := MakeAbs(p)
		t.Logf("Trial %d: %s -> %s", i, p, abs)
	}
	for i, p := range trials {
		abs, e := MakeAbsIfPathlike(p)
		if e != nil && e == ErrNotPathlike {
			t.Logf("Trial %d: %s -> %s (ErrNotPathlike)", i, p, abs)
		} else {
			t.Logf("Trial %d: %s -> %s", i, p, abs)
		}
	}
}

func TestPathIdentification(t *testing.T) {
	trials := []string{
		"``", `\`, "''", "_", "=",
	}
	_ = trials
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
