package config

import "testing"

func testInput() map[string]OptionKey {
	return map[string]OptionKey{
		"copyalldir":         BoolCopyAllDirs,
		"ignorehidden":       BoolIgnoreHidden,
		"useglobal":          BoolUseGlobalTarget,
		"baddabingbaddaboom": NotAnOption,
		"globaltarget":       StringGlobalTargetPath, // hits LookupExacts
		"nohiddenrepo":       NotAnOption,            // ambiguous: IgnoreHidden + IgnoreRepo
		"useglobaltgtdir":    NotAnOption,            // ambiguous: ForceGlobalTarget + GlobalTargetPath
		"override":           BoolOverrideOn,
		"killglobal":         BoolKillGlobalTarget,
		"killglobaltarget":   BoolKillGlobalTarget,
	}
}

func TestLookup(t *testing.T) {
	resultmap := testInput()
	for s, oKey := range resultmap {
		oLookup := LookupOption(s)
		t.Logf("lookup '%s'-> got %v", s, oLookup)
		if oLookup != oKey {
			t.Errorf("Wrong Val (got %s, wanted %s)", oLookup.String(), oKey.String())
		}
	}
}

func TestLookupSelf(t *testing.T) {
	for key, opt := range AllOptions {
		olookup := LookupOption(opt.NameText)
		t.Logf("search %s, got %s", opt.NameText, olookup.String())
		if olookup != key {
			t.Errorf("Wrong - got %s, wanted %s", olookup.String(), key.String())
		}
	}
}

func TestAllOptions(t *testing.T) {
	nfails := 0
	for i, k := range AllOptionIDs() {
		opt, ok := AllOptions[k]
		if !ok {
			t.Errorf("FAILURE: Option [%v] NOT A KEY IN AllOptions", k)
			nfails++

		} else {
			t.Logf("AllOption[%d]=%d;  name = %s", i, k, opt.NameText)
		}
	}
	if nfails > 0 {
		t.Logf("%d FAILING OptionKeys", nfails)
	}
}
