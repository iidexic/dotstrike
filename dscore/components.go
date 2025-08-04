package dscore

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	pops "iidexic.dotstrike/pathops"
)

// Denote whether paths in pathObjects are path or dir
type pathType byte

// Denote whether component is source or target. Uncertain of implementation
type componentType byte

const (
	_ pathType = iota
	filePath
	dirPath
	dnePath
)

const (
	_ componentType = iota
	specComponent
	sourceComponent
	targetComponent
	//linkedSource
	//linkedTarget
	//GlobalTarget
)

func (p pathType) string() string {
	switch p {
	case filePath:
		return "filepath"
	case dirPath:
		return "dirpath"
	case dnePath:
		return "path-DNE"
	}
	return ("unknown")
}
func (c componentType) string() string {
	switch c {
	case specComponent:
		return "spec"
	case sourceComponent:
		return "source"
	case targetComponent:
		return "target"
	}
	return "not a component type"
}

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
	Ptype   pathType      `toml:"ptype"`   //targetComponent requires dirPath //TODO: init Ptype or remove if not useful
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

func (pc pathComponent) Detail() string {
	lines := make([]string, 0, 16)
	//ctype := pc.Ctype.string()
	//header := fmt.Sprintf("Component: %s", ctype)
	path := fmt.Sprintf("	Path: %s (path type: %s)", pc.Path, pc.Ptype.string())
	parent := "	Parent Alias = " + pc.Parent
	lines = append(lines /*, header*/, path, parent)
	if pc.Alias != "" {
		lines = append(lines, fmt.Sprintf("	Alias: '%s'", pc.Alias))
	}
	if iqty := len(pc.Ignores); iqty > 0 {
		ig := make([]string, 0, iqty)
		ig = append(ig, "	List ignore patterns:")
		for i, pat := range pc.Ignores {
			ig = append(ig, fmt.Sprintf("	 [%d]: '%s'", i, pat))
		}
		lines = append(lines, ig...)
	}
	return strings.Join(lines, "\n")
}

// // id makes pathComponent Alias based on parent, ctype, path
// func (pc pathComponent) id() string {
// 	return pc.Parent + "" + string(pc.Ctype) + "" + pc.Alias
// }

// MatchesID determines whether the provided identifier string matches any of the following:
//   - Path
//   - AbsPath
//   - Alias
//   - BaseName of Abspath
func (pc pathComponent) MatchesID(checkid string) bool {
	return checkid == pc.Abspath || checkid == pc.Path || checkid == pc.Alias || checkid == pops.BaseName(pc.Abspath)
}
func (pc pathComponent) MatchesPath(id string) bool     { return id == pc.Abspath || id == pc.Path }
func (pc pathComponent) MatchesAlias(id string) bool    { return standardizeAlias(id) == pc.Alias }
func (pc pathComponent) MatchesPathBase(id string) bool { return id == pops.BaseName(pc.Abspath) }

func (pc pathComponent) IsSource() bool { return pc.Ctype == sourceComponent }

// ── Equality Check ──────────────────────────────────────────────────

func pathComponentEqual(pc, pc2 pathComponent) bool {
	return pc.Alias == pc2.Alias && pc.Abspath == pc2.Abspath && pc.Path == pc2.Path &&
		pc.Ptype == pc2.Ptype && pc.Ctype == pc2.Ctype && slices.Equal(pc.Ignores, pc2.Ignores)
}

// specEqual compares two cfg params for equality.
// standalone function to ensure compatible with slices.EqualFunc
func specEqual(S, S2 Spec) bool {
	return S.Alias == S2.Alias && S.Overrides == S2.Overrides && S.Ctype == S2.Ctype &&
		slices.EqualFunc(S.Sources, S2.Sources, pathComponentEqual) &&
		slices.EqualFunc(S.Targets, S2.Targets, pathComponentEqual) &&
		slices.Equal(S.Ignorepat, S2.Ignorepat)
}
