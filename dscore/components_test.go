package dscore

import (
	"testing"

	pops "iidexic.dotstrike/pathops"
)

var tcfg = cfg{
	Alias:   "testfiles",
	Sources: []pathComponent{{Path: "C:/dev/github/testfiles-in", Ctype: sourceComponent}},
	Targets: []pathComponent{{Path: "C:/dev/github/testfiles-out"}},
}

var cftest cfg = cfg{
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
---CFG---
Alias: %s, Ctype: %v
overrides: %+v`, tcfg.Alias, tcfg.Ctype, tcfg.Overrides)
	s := tcfg.Sources[0]
	t.Logf(`---Source---
Alias: %s, Abspath: %s, Path: %s
Ignores: %v
Ptype: %v,Ctype: %v`, s.Alias, s.Abspath, s.Path, s.Ignores, s.Ptype, s.Ctype)
	tcfg.initializeInherent()
}
