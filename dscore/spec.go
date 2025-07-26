package dscore

import (
	"fmt"
	"strings"

	pops "iidexic.dotstrike/pathops"
)

// TODO: What is this. Make this an error or get rid of it
func recoverReturn(explainer string) bool {
	if r := recover(); r != nil {
		fmt.Printf("FAILED [%s]. ERROR", explainer)
		fmt.Println(r)
		return true
	}
	return false
}

// Primary user data structure; contains info required to perform a dircopy operation
type Spec struct {
	Alias      string          `toml:"alias"`      // name, unique
	Sources    []pathComponent `toml:"sources"`    // paths marked as origin points
	Targets    []pathComponent `toml:"targets"`    // paths  marked as destination points
	Ignorepat  []string        `toml:"ignores"`    // ignorepat that apply to all sources
	OverrideOn bool            `toml:"overrideOn"` // enable overrides, prevent Overrides being over-written
	Overrides  prefs           `toml:"overrides"`  // override global prefs
	Ctype      componentType
}

// initializeInherent attributes of spec and child pathComponents
func (S *Spec) initializeInherent() {
	S.Ctype = specComponent
	for i := range S.Sources {
		S.Sources[i].Parent = S.Alias
		S.Sources[i].Ctype = sourceComponent
	}
	for i := range S.Targets {
		S.Targets[i].Parent = S.Alias
		S.Targets[i].Ctype = targetComponent
	}
}

// allInitialized check all source/target components to ensure all are initialized
func (S Spec) allInitialized() bool {
	all := S.Ctype > 0
	for _, src := range S.Sources {
		all = all && src.isInitialized()
	}
	for _, tgt := range S.Sources {
		all = all && tgt.isInitialized()
	}
	return all
}

// ── Find/Get Spec Info ──────────────────────────────────────────────

func (S Spec) getAlias() string        { return S.Alias }
func (S Spec) getCtype() componentType { return S.Ctype }

func (S *Spec) getSource(alias string) *pathComponent {
	for _, src := range S.Sources {
		if alias == src.Alias {
			return &src
		}
	}
	return nil
}
func (S *Spec) Detail() string {
	lines := make([]string, 0, 32)
	lines = append(lines, "Spec: "+S.Alias, "-------------------")

	// ── Sources ──────────────────────────────
	for _, src := range S.Sources {
		lines = append(lines, src.Detail())
	}

	// ── Targets ──────────────────────────────
	for _, tgt := range S.Targets {
		lines = append(lines, tgt.Detail())
	}

	// ── overrides ────────────────────────────
	overrideOn := fmt.Sprintf("Overrides Enabled: %t", S.OverrideOn)
	if !S.Overrides.equal(gd.data.Prefs) {
		lines = append(lines, fmt.Sprintf(`%s
Overrides:
	Keep Repo: %t
	Keep Hidden Files: %t
	Use Global Target: %t`, overrideOn, S.Overrides.KeepRepo, S.Overrides.KeepHidden, S.Overrides.GlobalTarget))
	}

	// ── Ignores ──────────────────────────────
	lines = append(lines, "Ignore Patterns:")
	for _, pat := range S.Ignorepat {
		lines = append(lines, "	- "+pat)
	}
	return strings.Join(lines, "\n")
}

func (S *Spec) getTarget(alias string) *pathComponent {
	for _, tgt := range S.Targets {
		if alias == tgt.Alias {
			return &tgt
		}
	}
	return nil
}
func (S *Spec) GetIgnores() *[]string { return &S.Ignorepat }

func (S *Spec) GetLocalPrefs() *prefs { return &S.Overrides }

// TODO: implement (in globalModify.go)

// IsPathChild looks for the path within the Spec's pathComponent slices
func (S Spec) IsPathChild(path string) bool {
	for _, src := range S.Sources {
		if src.Alias == path || src.Path == pops.MakeAbs(path) {
			return true
		}
	}
	for _, tgt := range S.Targets {
		if tgt.Alias == path || tgt.Path == pops.MakeAbs(path) {
			return true
		}
	}
	return false
}

// GetIfChild returns a pointer to the child source or target with the path or alias passed. Returns nil if none found
func (S Spec) GetIfChild(identifier string) *pathComponent {
	for _, src := range S.Sources {
		if src.Alias == identifier || src.Path == pops.MakeAbs(identifier) {
			return &src
		}
	}
	for _, tgt := range S.Targets {
		if tgt.Alias == identifier || tgt.Path == pops.MakeAbs(identifier) {
			return &tgt
		}
	}
	return nil
}

// ── Modifying Spec Data ─────────────────────────────────────────────

func (S *Spec) AddIgnores(ignores []string) {
	S.Ignorepat = append(S.Ignorepat, ignores...)
}

func (S *Spec) CheckAddPath(path string, isSource bool) bool {
	if !S.IsPathChild(path) {
		if isSource {
			S.Sources = append(S.Sources, *newPathComponent(path, sourceComponent))
		} else {
			S.Targets = append(S.Targets, *newPathComponent(path, targetComponent))
		}
	}
	return false
}

// CheckAddMultiplePaths adds paths to spec.Sources if isSource, spec.Targets if !isSource
func (S *Spec) CheckAddMultiplePaths(paths []string, isSource bool) {
	for _, p := range paths {
		S.CheckAddPath(p, isSource)
	}

}

// ── Running spec copy jobs ──────────────────────────────────────────
func (S Spec) RunCopy() error {
	if !S.allInitialized() {
		return fmt.Errorf("spec not initialized: %s", S.Alias)
	}
	copymachine := pops.GetCopierMaschine()
	//NOTE:
	for y, tgt := range S.Targets {
		for x, src := range S.Sources {
			//TODO: Finish Spec run copy jobs
			copymachine.NewJob(S.Alias+"."+fmt.Sprintf("%d", x)+"."+fmt.Sprintf("%d", y), src.Path, tgt.Path)
		}
	}
	return nil
}
