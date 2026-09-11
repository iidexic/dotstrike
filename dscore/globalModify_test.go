package dscore

import (
	"testing"
)

// ┌─────────────────────────────────────────────────────────┐
// │                          Tests                          │
// └─────────────────────────────────────────────────────────┘

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

