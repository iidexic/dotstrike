package dscore

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	pops "iidexic.dotstrike/pathops"
)

// Denote whether component is source or target. Uncertain of implementation
type componentType byte

const (
	_ componentType = iota
	specComponent
	sourceComponent
	targetComponent
	//linkedSource
	//linkedTarget
	//GlobalTarget
)

func (c componentType) String() string {
	switch c {
	case specComponent:
		return "spec"
	case sourceComponent:
		return "source"
	case targetComponent:
		return "target"
	}
	return "unknown"
}

var ErrComponentNotInitialized error = errors.New("Component not initialized")

// component interface. Only unused thing getting kept in
// type component interface {
// 	getAlias() string
// 	getCtype() componentType
// }

// PathComponent is the core of a source or target;
// contains path info
// TODO: Refactor Source and Target; eliminate complexity
//   - there are a lot of conditional situations here already
//   - there are also a lot of implementation inconsistencies.
type PathComponent struct {
	Alias   string        `toml:"alias"`
	Abspath string        `toml:"abspath"`
	Path    string        `toml:"path"`
	Ignores preIgnoreList `toml:"ignores"`
	Ctype   componentType `toml:"ctype"` //NOTE: Inherent? it is used in a couple places, replace with type assertions maybe
	Parent  string        //NOTE: INITIALIZE INHERENT
}

// isInitialized to check pc inherent-initialized. This is performed during startup and should never be false
func (pc PathComponent) isInitialized() bool { return pc.Parent != "" && pc.Ctype > 0 }

// interface methods
func (pc PathComponent) getAlias() string        { return pc.Alias }
func (pc PathComponent) getCtype() componentType { return pc.Ctype }

// TODO: replace this or replace the MakeAbs call
func newPathComponent(ospath string, ctype componentType) *PathComponent {
	apath := pops.MakeAbs(ospath)
	return &PathComponent{Path: apath, Ctype: ctype}
}

func (pc PathComponent) Detail() string {
	lines := make([]string, 0, 16)
	//ctype := pc.Ctype.string()
	//header := fmt.Sprintf("Component: %s", ctype)
	path := fmt.Sprintf("	Path: %s", pc.Path)
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

// TODO add basepath bool arg (huh?)

// MatchesID determines whether the provided identifier string matches any of the following:
//   - Path
//   - AbsPath
//   - Alias
//   - BaseName of Abspath
func (pc PathComponent) MatchesID(id string) bool {
	if pc.Abspath == "" {
		return id == pc.Abspath || id == pc.Path || id == pc.Alias ||
			strings.ToLower(id) == pops.Base(pc.Path)
	} else {
		return id == pc.Abspath || id == pc.Path || id == pc.Alias ||
			strings.ToLower(id) == pops.Base(pc.Abspath)
	}

}

func (pc PathComponent) IsSource() bool { return pc.Ctype == sourceComponent }

// ── Equality Check ──────────────────────────────────────────────────

func pathComponentEqual(pc, pc2 PathComponent) bool {
	return pc.Alias == pc2.Alias && pc.Abspath == pc2.Abspath && pc.Path == pc2.Path &&
		pc.Ctype == pc2.Ctype && slices.Equal(pc.Ignores, pc2.Ignores)
}

// specEqual compares two spec params for equality.
// standalone function to ensure compatible with slices.EqualFunc
func specEqual(S, S2 Spec) bool {
	return S.Alias == S2.Alias && S.Ctype == S2.Ctype &&
		slices.EqualFunc(S.Sources, S2.Sources, pathComponentEqual) &&
		slices.EqualFunc(S.Targets, S2.Targets, pathComponentEqual) &&
		slices.Equal(S.Ignorepat, S2.Ignorepat) &&
		S.Overrides.equal(S2.Overrides)
}
