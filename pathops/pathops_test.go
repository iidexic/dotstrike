package pops

import (
	"path/filepath"
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
