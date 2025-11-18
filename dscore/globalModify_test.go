package dscore

import (
	"testing"

	pops "iidexic.dotstrike/pathops"
)

var testTOMLpath = `D:\coding\github\dotstrike\_xtra\[samplefiles]\test_dotstrikeData.toml`

// ┌─────────────────────────────────────────────────────────┐
// │                          Tests                          │
// └─────────────────────────────────────────────────────────┘

// test making new spec and adding details. Unnecessary as is covered by TestEditEncode
func TestNewSpec(t *testing.T) {
	temp := initForTest(t)
	snew, e := temp.NewSpecEmpty("testnew")
	if e != nil {
		t.Logf("NewSpec Error: %s", e.Error())
	}
	if snew == nil {
		t.Error("spec returned nil")
	}

	snew.CheckAddPath("C:/Bingo.com/", true)
	snew.AddIgnores([]string{"yo", "*Bingo*"})
	t.Logf("temp spec list: %v", temp.Specs)
	t.Logf("tempData spec list: %v", tempData.Specs)

}

func TestEncodeHardAssign(t *testing.T) {
	p := gd.data.Prefs
	t.Log("globaldefaults are:")
	t.Log("specs empty, selected = 0")
	t.Logf("prefs:\n%+v", p)
	InitTempData()
	specNvim := Spec{
		Alias: "nvim", Sources: []PathComponent{{Path: "C:\\users\\derek\\appdata\\local\\nvim"}},
		Targets: []PathComponent{{Path: "@GLOBAL@"}}}
	specNvim.initializeInherent()
	specNvim.Sources[0].Alias = "nvim-config"
	// if !specNvim.allInitialized() {
	// 	t.Errorf("Spec '%s' not initialized:\n%+v", specNvim.Alias, specNvim)
	// }

	specWezterm := Spec{
		Alias: "wezterm", Sources: []PathComponent{{Path: "~\\.config\\wezterm"}},
		Targets: []PathComponent{{Path: "@GLOBAL@"}}}

	specWezterm.initializeInherent()
	// if !specWezterm.allInitialized() {
	// 	t.Errorf("Spec '%s' not initialized:\n%+v", specNvim.Alias, specNvim)
	// }
	tempData.Specs = append(tempData.Specs, specNvim)
	tempData.Specs = append(tempData.Specs, specWezterm)
	e := encodeTestfile(testTOMLpath, tempData.globalData)
	if e != nil {
		t.Errorf("Encode Error:%v", e)
	}
}

func TestGlobalEncodeSoftAssign(t *testing.T) {
	//LoadGlobals() // need to run CoreConfig?
	InitTempData()
	t.Log("Performed Init")
	tmp := TempData()
	st1, err := tmp.NewSpec("gamer", []string{"C:\\users\\derek\\appdata\\local\\nvim"}, []string{})
	if err != nil {
		t.Error(err)
	}

	failed := st1.Overrides.setOptMap(map[string]bool{"useglobaltarget": true})
	if len(failed) > 0 {
		t.Error("failed set option globaltarget")
	}
	if tmp.GetSpec("gamer") == nil {
		t.Error("nil pointer from created spec")
	}
	if !tmp.Modified {
		t.Error("Fail: TempData not marked as modified")
	}

}

func TestPrefSetByName(t *testing.T) {
	LoadGlobals()
	InitTempData()
	set := map[string]bool{"ignorehidden": false, "nohidden": true, "copyalldirs": true}
	spec := tempData.SelectedSpec()
	for k, v := range set {
		e := spec.Overrides.setByName(k, v)
		if e != nil {
			t.Errorf("failed setting %s = %t (%v)\nError:%s", k, v, v, e.Error())
		}
	}
}

func TestOptionID(t *testing.T) {
	tnames := []string{"ignorehidden", "nohidden", "copyalldirs"}
	expect := []ConfigOption{BoolIgnoreHidden, BoolIgnoreHidden, BoolCopyAllDirs}

	for i, nm := range tnames {
		found := OptionID(nm)
		if found != expect[i] {
			t.Errorf("expecting %s, found %s", expect[i].String(), found.String())
		}
	}
}

func TestSetOverridesMap(t *testing.T) {
	LoadGlobals()
	InitTempData()
	if !tempData.initialized {
		t.Errorf("tempData not initialized")
	}
	temp := TempData()
	if !temp.initialized {
		t.Errorf("TempData() not initialized")
	}
	set := map[string]bool{"ignorehidden": false, "nohidden": true}

	spec := temp.SelectedSpec()
	if spec != nil {
		setmap := tempData.SetSpecOverridesMap(spec, set)
		if len(setmap) > 0 {
			t.Errorf("setmap returned: %v", setmap)
		}
	} else {
		t.Errorf("SelectedSpec is nil")
	}

}

// test edit and encode; encodes to buffer and prints before manually writing to file
// WARN: Test Uncertain; not tested since major changes
func TestEncodeToBuffer(t *testing.T) {
	temp := initForTest(t)
	snew, e := temp.NewSpecEmpty("testEditEncodeSpec")
	if e != nil {
		t.Logf("NewSpec Errored: %s", e.Error())
	}
	if snew == nil {
		t.Error("globalModify.NewSpec() returned nil")
	}
	t.Log("Current State")
	t.Logf("%+v", temp)
	t.Logf("temp.Specs = %+v", temp.Specs)
	//Checking new spec and find spec
	testGet := temp.GetSpec("testEditEncodeSpec")
	if testGet == nil {
		t.Error("did not get ptr to spec")
	}

	testGet.AddIgnores([]string{"EE", "DD"})
	if !specEqual(temp.Specs[len(temp.Specs)-1], *testGet) {
		t.Errorf("temp.Spec '%s' != testSpec", temp.Specs[len(temp.Specs)-1].Alias)
		//do pointers match:
		if &temp.Specs[len(temp.Specs)-1] != testGet {
			t.Logf("Pointers don't match: %v != %v", &temp.Specs[len(temp.Specs)-1], testGet)
		}
		t.Logf("Last in temp.Specs:\n%+v", temp.Specs[len(temp.Specs)-1])
		t.Logf("Spec from GetModifiable:\n%+v", *testGet)
		for _, s := range temp.Specs {
			t.Logf("SPEC: %s, Ignores:%v", s.Alias, s.Ignorepat)
		}
	}
	buf, e := encodeToBuffer(temp.globalData)
	if e != nil {
		t.Errorf("encode error\n %s", e.Error())
	}
	t.Log("============[ FINAL BUFFER ]==============")
	t.Log(buf.String())
	f, e := pops.OpenFileRW(testTOMLpath)
	if e != nil {
		t.Errorf("file open error\n %s", e.Error())
	}
	defer f.Close()
	f.Write(buf.Bytes())
}

func TestEditEncode(t *testing.T) {
	temp := initForTest(t)
	snew, e := temp.NewSpecEmpty("testEditEncodeSpec")
	if e != nil {
		t.Logf("NewSpec Errored: %s", e.Error())
	}
	if snew == nil {
		t.Error("globalModify.NewSpec() returned nil")
	}
	t.Log("Current State of temp:")
	t.Logf("%+v", temp)
	t.Logf("temp.Specs = %+v", temp.Specs)
	testGet := temp.GetSpec("testEditEncodeSpec")
	if testGet == nil {
		t.Error("did not get ptr to spec")
	}

	testGet.AddIgnores([]string{"EE", "DD"})
	if !specEqual(temp.Specs[len(temp.Specs)-1], *testGet) {
		t.Errorf("temp.Spec '%s' != testSpec", temp.Specs[len(temp.Specs)-1].Alias)
		t.Logf("Last in temp.Specs:\n%+v", temp.Specs[len(temp.Specs)-1])
		t.Logf("Spec from GetModifiable:\n%+v", *testGet)
	}
	//e := gd.EncodeIfNeeded(temp)
	// if e != nil {
	// 	t.Errorf("Encode Error: %v", e)
	// }
	e = nil
	e = encodeTomltesting(testTOMLpath, temp.globalData)
	if e != nil {
		t.Errorf("encode error\n %s", e.Error())
	}
	//NOTE: Regardless, we need to set the original back, or we need to be using a test file

}
