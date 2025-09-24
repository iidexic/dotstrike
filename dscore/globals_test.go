package dscore

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	pops "iidexic.dotstrike/pathops"
)

var (
	errorEmpty  = fmt.Errorf("toml file exists but has no data")
	errorNoToml = fmt.Errorf("no toml file at path")
)

// same as globalModify_test.go encodeTomltesting
func encodeTestfile(path string, data *globalData) error {
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

func TestLoadOrEncodeDefaults(t *testing.T) {
	e := loadTestconfig(t)
	if e != nil {
		t.Logf("failed decoding datafile. Encoding defaults")
		encodeDefaultsToTestfile(t)
	} else {
		t.Log("Decoded Succesfully. result:")
		t.Logf("%+v", gd.data)
	}

}

func TestCoreConfig(t *testing.T) {
	t.Log(gd.Detail())
	CoreConfig()
	t.Log(gd.Detail())
}
func TestForceEncodeDefaults(t *testing.T) {
	encodeDefaultsToTestfile(t)
}

func loadTestconfig(t *testing.T) error {
	// Executs from dscore subdir, so using ../ below
	abstestdir := pops.MakeAbs("../_xtra/[samplefiles]")
	gotConfig := gd.GetConfig(abstestdir) //NOTE: uses dotstrikeData.toml

	fpath := pops.Joinpath(abstestdir, globalsFilename)
	t.Logf("filepath used in GetConfig: '%s'", fpath)
	if gotConfig {
		gd.status = badToml
		gd.decodeRawData()
		gd.loaded = true
		for _, c := range gd.data.Specs {
			c.initializeInherent()
		}
		undecoded := gd.md.Undecoded()
		if len(undecoded) > 0 {
			// real function logs into global struct
			t.Logf("undecoded values from .toml:")
			for i, u := range undecoded {
				t.Logf("%d) %s", i, u.String())
			}
		}
		if len(undecoded) < len(gd.md.Keys()) {
			// weird definition of success but ok
			gd.status = success
			return nil
		} else if len(gd.md.Keys()) == 0 {
			tLogErr(fmt.Sprintf("Path %s", fpath), errorEmpty, t)
			return errorEmpty
		}

	}
	tLogErr(fmt.Sprintf("Path %s", fpath), errorNoToml, t)
	return errorNoToml
}

func encodeDefaultsToTestfile(t *testing.T) {
	ee := encodeTestfile(testTOMLpath, &gd.data)
	if ee != nil {
		t.Errorf(`Failed writing default config to test toml.
Test data file could not be found/accecssed/created.`)
		t.Logf("ERROR: %v\n%s", ee, ee.Error())
	}
}
