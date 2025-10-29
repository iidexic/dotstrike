package dscore

import (
	"fmt"
	"path"

	pops "iidexic.dotstrike/pathops"
)

// enum type
type globalsReadResult int

// potential outcomes of attempting to read and load global config/user data into usable components
const ( // TODO:(hi) Change globalsReadResult into an error struct. Make these vars

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
			Bools: map[ConfigOption]bool{
				BoolIgnoreHidden:    false,
				BoolIgnoreRepo:      false,
				BoolUseGlobalTarget: true,
				//BoolSeparateSources: true,
				BoolCopyAllDirs: false,
				BoolNoFiles:     false,
			},
		},
		GlobalTargetPath: "~\\dotstrike\\globalTarget\\", // this doesnt work until transformed in CoreConfig.
		Specs:            []Spec{},
	},
}

// globalsFilename is the file that ds looks to pull settings and userdata from
const globalsFilename = "dotstrikeData.toml"
const globalDirHomeRelative = ".config/dotstrike"
const globalDirConfigRelative = "/dotstrike"

var GlobalConfigPath string

func globalsFilepath() string {
	if GlobalConfigPath != "" {
		return GlobalConfigPath
	}
	if gd.dsconfigPath != "" {
		return gd.dsconfigPath
	}
	panic(fmt.Errorf("Global config path not set"))
}

func ConfigTomlPath() string {
	if gd.loaded && gd.dsconfigPath != "" {
		return gd.dsconfigPath
	} else if !gd.loaded {
		return "CONFIG NOT LOADED"
	} else {
		return "CONFIG PATH NOT POPULATED"
	}
}

func (g *globalData) getSpec(alias string) *Spec {
	for _, s := range g.Specs {
		if s.Alias == alias {
			return &s
		}
	}
	return nil
}
func (G *globals) GetConfigFrom(filepath string) bool {
	readfile := pops.ReadFile(filepath)
	if !readfile.Failed() {
		G.rawContents = string(readfile.Contents)
		G.dsconfigPath = filepath
		return true
	} else if readfile.Fail == pops.FailedOpen {
		G.status = badRead
	} else if readfile.Fail == pops.FileNotExist {
		G.status = noInit
	}
	G.logG(readfile.Fail.Detail())
	return false
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

// TODO: (low) rewrite

// CoreConfig called to find ds data file in all possible locations
// ?TODO: If no config file exists, create one and encode gd defaults
func (G *globals) makeCfgPath(suffix string) string {
	if !pops.HaveHome() && pops.ErrGetHomedir == nil && pops.ErrGetConfigdir == nil {
		pops.GetSysDirs()
	}
	return pops.HomeJoinC(suffix)
}

func LoadGlobals() error {
	var errinit error
	INIT := initializer{
		filename:      globalsFilename,
		SysFileErrors: make(map[string]error),
	}
	e := INIT.Config()
	if e != nil {
		errinit = e
		if e == pops.ErrorUserDirs {
			gd.logfG("System config/home directories not found in env")
		} else {
			gd.logfG("Failed to load config: %s", e.Error())
			e = INIT.WriteConfigDefaults()
			if e != nil {
				errinit = extenderror(errinit, e, "write defaults to toml")
				gd.logfG("%s", e.Error())
			} else {
				gd.logfG("Wrote defaults to toml at %s", GlobalConfigPath)
			}
		}
	}
	return e
}

func CoreConfig() error {
	if pops.HomePath == nil {
		pops.GetSysDirs()
	}

	//TODO: clean up this homepath/GlobalTargetPath solution
	cfgdir := gd.makeCfgPath(globalDirHomeRelative)
	// for default config
	gd.data.GlobalTargetPath = pops.TildeExpand(gd.data.GlobalTargetPath)
	gotConfig := gd.GetConfig(cfgdir)
	if gotConfig {
		gd.status = badToml //pre-emptive
		gd.decodeRawData()
		gd.loaded = true
		// better way to do this?
		for _, c := range gd.data.Specs {
			c.initializeInherent() // BUG:this isnt gonna work duh
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
	return nil

}

func EndEncode() {
	if tempData.Modified {
		e := tempData.encodeModified()
		if e != nil {
			if pe, ok := any(e).(pops.PathError); ok {
				gd.forcelogG(fmt.Sprintf("Error opening file (OpenRW): %v", pe))
			} else {

				gd.forcelogG(fmt.Sprintf("Error encoding to TOML file: %s", e.Error()))

			}
		}
	}
}
