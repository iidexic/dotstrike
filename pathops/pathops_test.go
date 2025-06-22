package pops

import (
	"path/filepath"
	"slices"
	"testing"
)

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
