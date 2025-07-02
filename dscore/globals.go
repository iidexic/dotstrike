package dscore

import (
	"fmt"
	"os"
	"path"

	pops "iidexic.dotstrike/pathops"
)

// enum type
type globalsReadResult int

// TODO: Just turn these into errors. Basically are, just called it .string() instead of .Error()
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
var gd = globals{
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
const globalPathHomeRelative = ".config/dotstrike/dotstrikeData.toml"

func globalsFilepath() string {
	gpath, e := pops.HomeJoin(globalPathHomeRelative)
	if e != nil {
		panic(e)
	}
	return gpath
}

func GetGlobals() (*globals, error) {
	if gd.loaded {
		return &gd, nil
	}
	return &globals{}, fmt.Errorf("Globals not loaded.\n Globals = %+v", gd)
}
func (g *globalData) GetCfg(alias string) *cfg {
	for _, cc := range g.Cfgs {
		if cc.Alias == alias {
			return &cc
		}
	}
	return nil
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
	gotConfig := gd.GetConfig(cfgdir)
	if gotConfig {
		gd.status = badToml //pre-emptive
		gd.decodeRawData()
		gd.loaded = true
		// better way to do this?
		for _, c := range gd.data.Cfgs {
			c.initializeInherent()
		}
		undecoded := gd.md.Undecoded()
		if len(undecoded) > 0 {
			gd.logfG("undecoded values from .toml:\n%+v", undecoded)
		}
		if len(undecoded) < len(gd.md.Keys()) {
			gd.status = success
		}

	} else {

	}

}

func EndEncode() {
}
