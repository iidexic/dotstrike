package dscore

import (
	"fmt"
)

// Denote whether paths in pathObjects are path or dir
type pathType int

// Denote whether component is source or target. Uncertain of implementation
type componentType int

const (
	filePath pathType = iota
	dirPath
	sourceComponent componentType = iota
	targetComponent
	cfgComponent
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

func (pc pathComponent) getAlias() string        { return pc.Alias }
func (pc pathComponent) getCtype() componentType { return pc.Ctype }

// cfg is the primary structure used to define a move/strike
type cfg struct {
	Alias     string          `toml:"alias"`     // name, unique
	Sources   []pathComponent `toml:"sources"`   // files or directories marked as origin points
	Targets   []pathComponent `toml:"targets"`   // files or directories marked as destination points
	Ignorepat []string        `toml:"ignores"`   // ignore patterns that apply to all sources
	Overrides prefs           `toml:"overrides"` //map of settings that will be prioritized over global set
	Ctype     componentType
}

func (cc cfg) getAlias() string        { return cc.Alias }
func (cc cfg) getCtype() componentType { return cc.Ctype }
func (cc cfg) status() string {
	expln := fmt.Sprintf("cfg:'%s' - Sources:\n%+v", cc.Alias, cc.Sources)
	return expln
}
