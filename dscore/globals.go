package dscore

import (
	"os"
	"path"

	pops "iidexic.dotstrike/pathops"
)

// enum type
type globalsReadResult int

// potential outcomes of attempting to read and load global config/user data into usable components
const (
	preReadGlobalsReadResult = iota
	noInit
	BadToml
	success
)

func (gr globalsReadResult) string() string {
	switch gr {
	case preReadGlobalsReadResult:
		return "PRE-READ: globals file read not triggered"
	case noInit:
		return "NO-INIT: Init file not found"
	case BadToml:
		return "BAD-CONFIG: globals/config file exists, but has malformed structure/syntax"
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
		Prefs: prefs{keepRepo: true, keepHidden: true, globalTarget: true},
		cfgs:  []cfg{}, CfgToml: []cfgMake{},
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

// GetConfig (from file) to be loaded into G
func (G *globals) GetConfig(dirpath string) bool {
	fpath := path.Join(dirpath, globalsFilename)
	//fread possible failureType (fread.Fail) = None or FileNotExist or FailedOpen
	fread := pops.ReadFile(fpath)
	// if ReadFile succeeded
	if fread.Fail == pops.None {
		if !G.loaded {
			G.rawContents = string(fread.Contents)
			G.dsconfigPath = fpath
			return true
		} else {
			G.checkedpaths = append(G.checkedpaths, dirpath)
		}

	} else if fread.Fail == pops.FailedOpen {
		G.output(fread.Fail.Detail())
	}
	return false
}

// CoreConfig called to find ds data file in all possible locations
// TODO: Update for MVP or remove
// might as well keep this
func CoreConfig() {
	homedir, errcfg := os.UserHomeDir()
	cfgdir := path.Join(homedir, ".config/dotstrike")
	// for now just panic here
	ifer(errcfg)
	gotConfig := GD.GetConfig(cfgdir)
	if gotConfig {
		GD.DecodeRawData()

	} else {
	}

}
