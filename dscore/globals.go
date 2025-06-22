package dscore

import (
	"fmt"
	"os"
	"path"

	pops "iidexic.dotstrike/pathops"
)

// enum type
type globalsReadResult int

// potential outcomes of attempting to read and load global config/user data into usable components
const (
	preInit = iota
	noInit
	badToml
	badRead
	extraInit
	success
)

func (gr globalsReadResult) string() string {
	switch gr {
	case preInit:
		return "PRE-INIT: globals file read not triggered"
	case noInit:
		return "NO-INIT: Init file not found"
	case badToml:
		return "BAD-CONFIG: globals/config file exists, but has malformed structure or syntax"
	case badRead:
		return fmt.Sprintf("BAD-READ: %s exists but read failed", globalsFilename)
	case success:
		return "SUCCESS"
	}
	return "UNKNOWN-CASE-ERROR"
}

// primary var that data will be pulled into.
// as is, it also serves as configuration defaults.
var GD = globals{
	status: noInit,
	loaded: false,
	data: globalData{
		Prefs: prefs{KeepRepo: true, KeepHidden: true, GlobalTarget: true},
		Cfgs:  []cfg{},
	},
}

// GD fields
/* status        globalsReadResult
loaded        bool
cfgs          []cfg
preferences   prefs
dsconfigPath  string
checkedpaths  []string
rawContents   []byte
GlobalMessage []string */

// globalsFilename is the file that ds looks to pull settings and userdata from
const globalsFilename = "dotstrikeData.toml"

func GetGlobals() (*globals, error) {
	if GD.loaded {
		return &GD, nil
	}
	return &globals{}, fmt.Errorf("Globals not loaded.\n Globals = %+v", GD)
}

// GetConfig reads dotstrikeData.toml in provided directory.
// on success: populates G.dsconfigPath, reads file into G.rawContents
func (G *globals) GetConfig(dirpath string) bool {
	fpath := path.Join(dirpath, globalsFilename)
	fread := pops.ReadFile(fpath) // None|FileNotExist|FailedOpen
	// if ReadFile succeeded
	if !fread.Failed() {
		if !G.loaded {
			G.rawContents = string(fread.Contents)
			G.dsconfigPath = fpath
			return true
		} else {
			G.checkedpaths = append(G.checkedpaths, dirpath)
			G.status = extraInit
			return false
		}

	} else if fread.Fail == pops.FailedOpen {
		G.status = badRead
	} else if fread.Fail == pops.FileNotExist {
		G.status = noInit
	}
	G.logG(fread.Fail.Detail())
	return false
}

// CoreConfig called to find ds data file in all possible locations
// TODO: Establish better separation of functionality with GetConfig
func CoreConfig() {
	homedir, errcfg := os.UserHomeDir()
	ifer(errcfg) // for now just panic
	cfgdir := path.Join(homedir, ".config/dotstrike")
	gotConfig := GD.GetConfig(cfgdir)
	if gotConfig {
		GD.status = badToml //pre-emptive
		GD.decodeRawData()
		GD.loaded = true
		undecoded := GD.md.Undecoded()
		if len(undecoded) > 0 {
			GD.logfG("undecoded values from .toml:\n%+v", undecoded)
		}
		if len(undecoded) < len(GD.md.Keys()) {
			GD.status = success
		}

	} else {

	}

}

func EndEncode() {
}
