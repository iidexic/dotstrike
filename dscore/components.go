package dscore

import (
	"fmt"
)

// Denote whether paths in pathObjects are path or dir
type pathType int

// Denote whether component is source or target. Uncertain of implementation
type componentType int

const (
	cfgComponent componentType = iota
	sourceComponent
	targetComponent
	filePath pathType = iota
	dirPath
)

// component interface for interop search
type component interface {
	getAlias() string
	getCtype() componentType
}

// pathComponent is the core of a source or target;
// contains path info
// TODO: Identify if separate structs are needed for source and target;
//   - there are a lot of conditional situations here already
type pathComponent struct {
	Path    string   `toml:"path"`
	Ptype   pathType `toml:"ptype"` //if targetComponent: required to be dirPath
	Ctype   componentType
	Abspath string `toml:"abspath"`
	// Ignores:
	// - if sourceComponent+dirPath: patterns(exact paths?) to ignore
	// - if targetComponent: patterns to avoid copying to this specific target
	Ignores []string `toml:"ignores"`
	Alias   string   `toml:"alias"`
}

// cfg is the primary structure used to define a move/strike
type cfg struct {
	Alias     string          `toml:"alias"`     // name, unique
	Sources   []pathComponent `toml:"sources"`   // files or directories marked as origin points
	Targets   []pathComponent `toml:"targets"`   // files or directories marked as destination points
	Ignorepat []string        `toml:"ignores"`   // ignore patterns that apply to all sources
	Overrides prefs           `toml:"overrides"` //map of settings; prioritized over global
	Ctype     componentType
}

func (pc pathComponent) getAlias() string        { return pc.Alias }
func (pc pathComponent) getCtype() componentType { return pc.Ctype }

func (cc cfg) getAlias() string        { return cc.Alias }
func (cc cfg) getCtype() componentType { return cc.Ctype }
func (cc cfg) status() string {
	expln := fmt.Sprintf("cfg:'%s' - Sources:\n%+v", cc.Alias, cc.Sources)
	return expln
}

type OpMode int

const (
	ModifyComponent OpMode = iota
	NewComponent
	DeleteComponent
)

type Operation struct {
	Get      Lookup //what need to find; can dictate optarget
	Mode     OpMode
	optarget componentType // the component type where direct modification is taking place
	modCfg   *cfg
	modPath  *pathComponent
}

func (O *Operation) ProcessFind() {

}
