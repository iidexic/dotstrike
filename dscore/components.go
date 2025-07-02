package dscore

import (
	"fmt"
	"slices"

	pops "iidexic.dotstrike/pathops"
)

// Denote whether paths in pathObjects are path or dir
type pathType = byte

// Denote whether component is source or target. Uncertain of implementation
type componentType = byte

const (
	cfgComponent componentType = iota
	sourceComponent
	targetComponent
	filePath pathType = iota
	dirPath
	dnePath
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
	Alias   string        `toml:"alias"`
	Abspath string        `toml:"abspath"`
	Path    string        `toml:"path"`
	Ignores []string      `toml:"ignores"` // remove if not using target ignore copy
	Ptype   pathType      `toml:"ptype"`   //targetComponent requires dirPath
	Ctype   componentType `toml:"ctype"`   //NOTE: not  implemented
	Parent  string        //NOTE: INITIALIZE INHERENT
}

// cfg is the primary structure used to define a move/strike
type cfg struct {
	Alias     string          `toml:"alias"`     // name, unique
	Sources   []pathComponent `toml:"sources"`   // paths marked as origin points
	Targets   []pathComponent `toml:"targets"`   // paths  marked as destination points
	Ignorepat []string        `toml:"ignores"`   // ignorepat that apply to all sources
	Overrides prefs           `toml:"overrides"` // override global prefs
	Ctype     componentType
}

// initializeInherent attributes of cfg and child pathComponents
func (cc *cfg) initializeInherent() {
	cc.Ctype = cfgComponent
	for _, src := range cc.Sources {
		src.Parent = cc.Alias
		src.Ptype = sourceComponent
	}
	for _, tgt := range cc.Targets {
		tgt.Parent = cc.Alias
		tgt.Ptype = targetComponent
	}
}

// interface methods
func (pc pathComponent) getAlias() string        { return pc.Alias }
func (pc pathComponent) getCtype() componentType { return pc.Ctype }
func (cc cfg) getAlias() string                  { return cc.Alias }
func (cc cfg) getCtype() componentType           { return cc.Ctype }

func (cc *cfg) getSource(alias string) *pathComponent {
	for _, src := range cc.Sources {
		if alias == src.Alias {
			return &src
		}
	}
	return nil
}
func (cc *cfg) getTarget(alias string) *pathComponent {
	for _, tgt := range cc.Targets {
		if alias == tgt.Alias {
			return &tgt
		}
	}
	return nil
}

func (cc cfg) status() string {
	expln := fmt.Sprintf("cfg:'%s' - Sources:\n%+v", cc.Alias, cc.Sources)
	return expln
}
func newPathComponent(ospath string, ctype componentType) *pathComponent {
	// do this here or at the end of modify?
	goodpath := pops.CleanPath(ospath)
	return &pathComponent{Path: goodpath, Ctype: ctype}

}

func (pc pathComponent) ID() string {
	return pc.Parent + "~>" + string(pc.Ctype) + "~>" + pc.Alias
}
func pathComponentEqual(pc, pc2 pathComponent) bool {
	return pc.Alias == pc2.Alias && pc.Abspath == pc2.Abspath && pc.Path == pc2.Path &&
		pc.Ptype == pc2.Ptype && pc.Ctype == pc2.Ctype && slices.Equal(pc.Ignores, pc2.Ignores)
}

// cfgEqual compares two cfg params for equality.
// standalone function to ensure compatible with slices.EqualFunc
func cfgEqual(cc, cc2 cfg) bool {
	return cc.Alias == cc2.Alias && cc.Overrides == cc2.Overrides && cc.Ctype == cc2.Ctype &&
		slices.EqualFunc(cc.Sources, cc2.Sources, pathComponentEqual) &&
		slices.EqualFunc(cc.Targets, cc2.Targets, pathComponentEqual) &&
		slices.Equal(cc.Ignorepat, cc2.Ignorepat)
}
