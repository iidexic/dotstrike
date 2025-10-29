package dscore

import (
	"testing"

	"iidexic.dotstrike/uout"
)

func dumpGlobalLog(t *testing.T) {
	out := uout.NewOut("[ Global Log ]")
	out.ILV(gd.GlobalMessage)
	t.Log(out.String())
}

func TestInitialize(t *testing.T) {
	I := initializer{
		filename:      globalsFilename,
		SysFileErrors: make(map[string]error),
	}
	e := I.Config()
	if e != nil {
		t.Log("THERE WAS AN ERROR THAT MADE IT BACK")
		t.Error(e)
		dumpGlobalLog(t)
		t.Logf("%v", gd)
	}

}

func TestLengthNamedPaths(t *testing.T) {
	cfgpaths := MakeSysConfigPaths(globalsFilename)
	t.Logf("len(cfgpaths): %d", len(cfgpaths))
	t.Logf("cfgpaths: %v", cfgpaths)
}

func TestConfigProcess(t *testing.T) {
	I := initializer{
		filename:      globalsFilename,
		SysFileErrors: make(map[string]error),
	}

	I.tomlpaths = MakeSysConfigPaths(globalsFilename)
	t.Logf("tomlpaths: %v", I.tomlpaths)

	e := I.populateGlobalData()
	if e != nil {
		t.Errorf("Init tomlpaths: %v", I.tomlpaths)
	}
	t.Logf("gd: %v", gd)
}
