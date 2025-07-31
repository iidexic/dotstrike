package dscore

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	pops "iidexic.dotstrike/pathops"
)

var testTOMLpath = `D:\coding\github\dotstrike\[samplefiles]\test_dotstrikeData.toml`

// ── Error logging ───────────────────────────────────────────────────
func tLogErr(context string, e error, t *testing.T) {
	if e != nil {
		t.Logf("%s: %s", context, e.Error())
	}
}

func tError(context string, e error, t *testing.T) {
	if e != nil {
		t.Errorf("%s: %s", context, e.Error())
	}
}

// ── Decoding ────────────────────────────────────────────────────────

func loadconfig(t *testing.T) *globalModify {
	CoreConfig()   // decode globals
	InitTempData() // load tempdata struct with globals details
	temp := TempData()
	if temp == nil {
		t.Error("Temp Data not initialized")
	}
	return temp
}

// ── Encoding ────────────────────────────────────────────────────────

func encodeTomltesting(path string, data *globalData) error {
	file, e := pops.OpenFileRW(pops.CleanPath(path))
	if file != nil {
		defer file.Close()
	}
	if e != nil || file == nil {
		return fmt.Errorf("Error opening toml for write: %w", e)
	}
	encode := toml.NewEncoder(file)
	e = encode.Encode(*data)
	if e != nil {
		return e
	} else {
		return nil
	}
}
func encodeToBuffer(data *globalData) (bytes.Buffer, error) {
	buf := bytes.Buffer{}
	e := toml.NewEncoder(&buf).Encode(*data)
	if e != nil {
		return buf, e
	}
	return buf, nil
}

// ┌─────────────────────────────────────────────────────────┐
// │                          Tests                          │
// └─────────────────────────────────────────────────────────┘

// test making new spec and adding details. Unnecessary as is covered by TestEditEncode
func TestNewSpec(t *testing.T) {
	temp := loadconfig(t)
	snew, e := temp.NewSpec("testnew")
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
		Alias: "nvim", Sources: []pathComponent{{Path: "C:\\users\\derek\\appdata\\local\\nvim"}},
		Targets: []pathComponent{{Path: "@GLOBAL@"}}}
	specNvim.initializeInherent()
	specNvim.Sources[0].Alias = "nvim-config"
	if !specNvim.allInitialized() {
		t.Errorf("Spec '%s' not initialized:\n%+v", specNvim.Alias, specNvim)
	}

	specWezterm := Spec{
		Alias: "wezterm", Sources: []pathComponent{{Path: "~\\.config\\wezterm"}},
		Targets: []pathComponent{{Path: "@GLOBAL@"}}}

	specWezterm.initializeInherent()
	if !specWezterm.allInitialized() {
		t.Errorf("Spec '%s' not initialized:\n%+v", specNvim.Alias, specNvim)
	}
	tempData.Specs = append(tempData.Specs, specNvim)
	tempData.Specs = append(tempData.Specs, specWezterm)
	encodeTestfile(testTOMLpath, tempData.globalData)
}

func TestGlobalEncodeSoftAssign(t *testing.T) {
	//CoreConfig() // need to run CoreConfig?
	InitTempData()
	t.Log("Performed Init")
	tmp := TempData()
	st1, err := tmp.NewSpec("gamer", "C:\\users\\derek\\appdata\\local\\nvim")
	if err != nil {
		t.Error(err)
	}

	err = st1.Overrides.SetM(map[string]bool{"globaltarget": true})
	if err != nil {
		t.Error(err)
	}
	if tmp.getSpec("gamer") == nil {
		t.Error("nil pointer from created spec")
	}
	if !tmp.Modified {
		t.Error("Fail: TempData not marked as modified")
	}

}

// test edit and encode; encodes to buffer and prints before manually writing to file
func TestEncodeToBuffer(t *testing.T) {
	temp := loadconfig(t)
	snew, e := temp.NewSpec("testEditEncodeSpec")
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
	testGet, eget := temp.GetModifiableSpec("testEditEncodeSpec")
	if eget != nil {
		t.Error(eget)
	} else if testGet == nil {
		t.Error("still not getting ptr to spec")
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
	temp := loadconfig(t)
	snew, e := temp.NewSpec("testEditEncodeSpec")
	if e != nil {
		t.Logf("NewSpec Errored: %s", e.Error())
	}
	if snew == nil {
		t.Error("globalModify.NewSpec() returned nil")
	}
	t.Log("Current State of temp:")
	t.Logf("%+v", temp)
	t.Logf("temp.Specs = %+v", temp.Specs)
	testGet, eget := temp.GetModifiableSpec("testEditEncodeSpec")
	if eget != nil {
		t.Error(eget)
	} else if testGet == nil {
		t.Error("still not getting ptr to spec")
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
