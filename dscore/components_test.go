package dscore

import (
	"testing"
)

var tspec = Spec{
	Alias:   "testfiles",
	Sources: []PathComponent{{Path: "C:/dev/github/testfiles-in", Ctype: sourceComponent}},
	Targets: []PathComponent{{Path: "C:/dev/github/testfiles-out"}},
}

func TestInherent(t *testing.T) {
	t.Logf(`component defaults:
---SPEC---
Alias: %s, Ctype: %v
overrides: %+v`, tspec.Alias, tspec.Ctype, tspec.Overrides)
	s := tspec.Sources[0]
	t.Logf(`---Source---
Alias: %s, Abspath: %s, Path: %s
Ignores: %v
Ctype: %v`, s.Alias, s.Abspath, s.Path, s.Ignores, s.Ctype)
	tspec.initializeInherent()
}
