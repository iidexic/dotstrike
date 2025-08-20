package config

import "testing"

func testInput() map[string]OptionKey {
	return map[string]OptionKey{
		"copydir":            BoolCopyAllDirs,
		"ignorehidden":       BoolIgnoreHidden,
		"useglobal":          BoolUseGlobalTarget,
		"baddabingbaddaboom": NotAnOption,
		"globaltarget":       NotAnOption, //TODO: fix this, will probably get annoying

		//BUG: Points to both IgnoreHidden and IgnoreRepo. Whichever is first  is it.
		//	Code defines IgnoreRepo first in the map
		//	EDIT: OUTCOME IS UNCERTAIN. MAY NEED AN IMPROVED LOOKUP
		"nohiddenrepo": BoolIgnoreRepo,
		//BUG: same bug as above (UseGlobalTarget/GlobalTargetPath)
		"useglobaltgtdir": BoolUseGlobalTarget,
		"override":        BoolOverrideOn,
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
