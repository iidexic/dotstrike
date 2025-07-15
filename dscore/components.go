package dscore

import (
	"errors"
	"slices"

	pops "iidexic.dotstrike/pathops"
)

// Denote whether paths in pathObjects are path or dir
type pathType = byte

// Denote whether component is source or target. Uncertain of implementation
type componentType = byte

const (
	_ pathType = iota
	filePath
	dirPath
	dnePath
)

const (
	_ componentType = iota
	cfgComponent
	sourceComponent
	targetComponent
	//linkedSource
	//linkedTarget
	//GlobalTarget
)

var ErrComponentNotInitialized error = errors.New("Component not initialized")

// component interface for interop search
type component interface {
	getAlias() string
	getCtype() componentType
}

// pathComponent is the core of a source or target;
// contains path info
// TODO: Refactor Source and Target; eliminate complexity
//   - there are a lot of conditional situations here already
//   - there are also a lot of implementation inconsistencies.
type pathComponent struct {
	Alias   string        `toml:"alias"`
	Abspath string        `toml:"abspath"`
	Path    string        `toml:"path"`
	Ignores []string      `toml:"ignores"` // remove if not using target ignore copy
	Ptype   pathType      `toml:"ptype"`   //targetComponent requires dirPath
	Ctype   componentType `toml:"ctype"`   //NOTE: not  implemented. Inherent??
	Parent  string        //NOTE: INITIALIZE INHERENT
}

// isInitialized to check pc inherent-initialized. This is performed during startup and should never be false
func (pc pathComponent) isInitialized() bool { return pc.Parent != "" && pc.Ctype > 0 }

// interface methods
func (pc pathComponent) getAlias() string        { return pc.Alias }
func (pc pathComponent) getCtype() componentType { return pc.Ctype }

func newPathComponent(ospath string, ctype componentType) *pathComponent {
	apath := pops.MakeAbs(ospath)
	return &pathComponent{Path: apath, Ctype: ctype}
}

// id makes pathComponent Alias based on parent, ctype, path
func (pc pathComponent) id() string {
	return pc.Parent + "" + string(pc.Ctype) + "" + pc.Alias
}

func (pc pathComponent) MatchesID(checkid string) bool {
	return checkid == pc.Abspath || checkid == pc.Path || checkid == pc.Alias || checkid == pops.BaseName(pc.Abspath)
}

// ── Equality Check ──────────────────────────────────────────────────

func pathComponentEqual(pc, pc2 pathComponent) bool {
	return pc.Alias == pc2.Alias && pc.Abspath == pc2.Abspath && pc.Path == pc2.Path &&
		pc.Ptype == pc2.Ptype && pc.Ctype == pc2.Ctype && slices.Equal(pc.Ignores, pc2.Ignores)
}

// specEqual compares two cfg params for equality.
// standalone function to ensure compatible with slices.EqualFunc
func specEqual(S, S2 spec) bool {
	return S.Alias == S2.Alias && S.Overrides == S2.Overrides && S.Ctype == S2.Ctype &&
		slices.EqualFunc(S.Sources, S2.Sources, pathComponentEqual) &&
		slices.EqualFunc(S.Targets, S2.Targets, pathComponentEqual) &&
		slices.Equal(S.Ignorepat, S2.Ignorepat)
}
