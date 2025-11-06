package dscore

import (
	"fmt"

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
				BoolUseGlobalTarget: false,
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

// TODO: # 2 0- FINISH
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

// EndEncode writes any changes made to tempData to the global toml file
//
// runs on cobra finalize
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
