package dscore

import (
	"testing"

	pops "iidexic.dotstrike/pathops"
)

var tspec = Spec{
	Alias:   "testfiles",
	Sources: []pathComponent{{Path: "C:/dev/github/testfiles-in", Ctype: sourceComponent}},
	Targets: []pathComponent{{Path: "C:/dev/github/testfiles-out"}},
}

var sftest Spec = Spec{
	Alias: "tilde_dotconfig",
	Sources: []pathComponent{
		{
			Alias:   "Source1",
			Path:    "~/.config/",
			Abspath: pops.HomeDirtyJoin(".config/"),
		},
	},
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
Ptype: %v,Ctype: %v`, s.Alias, s.Abspath, s.Path, s.Ignores, s.Ptype, s.Ctype)
	tspec.initializeInherent()
}
