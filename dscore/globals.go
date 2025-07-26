package dscore

import (
	"fmt"
	"path"

	pops "iidexic.dotstrike/pathops"
)

// enum type
type globalsReadResult int

// potential outcomes of attempting to read and load global config/user data into usable components
const ( // TODO: Change globalsReadResult into an error struct. Make these vars

	preInit = iota
	noInit
	badToml
	badRead
	extraInit
	success
)

func (gr globalsReadResult) Error() string {
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
		Prefs: prefs{
			KeepRepo:         true,
			KeepHidden:       true,
			GlobalTarget:     true,
			GlobalTargetPath: "~\\dotstrike\\globalTarget\\", // this doesnt work until transformed in CoreConfig.
		},
		Specs: []Spec{},
	},
}
var GetSpec = gd.data.getSpec

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
const globalDirHomeRelative = ".config/dotstrike"

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

func (g *globalData) getSpec(alias string) *Spec {
	for _, s := range g.Specs {
		if s.Alias == alias {
			return &s
		}
	}
	return nil
}

// TODO: move to (val, error) format so can directly diagnose os.ErrNotExist
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

// TODO: Establish better separation of functionality between CoreConfig and GetConfig

// CoreConfig called to find ds data file in all possible locations
// ?TODO: If no config file exists, create one and encode gd defaults
// TODO: require vars be passed (globalDir)?
func CoreConfig() {

	//TODO: clean up this homepath/GlobalTargetPath solution
	cfgdir := pops.Joinpath(*pops.HomePath, globalDirHomeRelative)
	gd.data.Prefs.GlobalTargetPath = pops.TildeDirty(gd.data.Prefs.GlobalTargetPath)
	gotConfig := gd.GetConfig(cfgdir)
	if gotConfig {
		gd.status = badToml //pre-emptive
		gd.decodeRawData()
		gd.loaded = true
		// better way to do this?
		for _, c := range gd.data.Specs {
			c.initializeInherent()
		}
		undecoded := gd.md.Undecoded()
		if len(undecoded) > 0 {
			gd.logfG("undecoded values from .toml:\n%+v", undecoded)
		}
		//TODO: better way of determining success
		if len(undecoded) < len(gd.md.Keys()) {
			gd.status = success
		}

	} else {
		// load modify and write defaults to file first thing
		// Option 1: use the function I wrote literally for this
		ee := gd.encodeDefaults()
		if ee != nil {
			panic(fmt.Errorf(
				`Failed writing default config to file (%w)
User data file not found and could not be made.`, ee))
		}

	}

}

func EndEncode() {
}
