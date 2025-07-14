package dscore

import (
	"fmt"

	pops "iidexic.dotstrike/pathops"
)

func recoverReturn(explainer string) bool {
	if r := recover(); r != nil {
		fmt.Printf("FAILED [%s]. ERROR", explainer)
		fmt.Println(r)
		return true
	}
	return false
}

type spec struct {
	Alias     string          `toml:"alias"`     // name, unique
	Sources   []pathComponent `toml:"sources"`   // paths marked as origin points
	Targets   []pathComponent `toml:"targets"`   // paths  marked as destination points
	Ignorepat []string        `toml:"ignores"`   // ignorepat that apply to all sources
	Overrides prefs           `toml:"overrides"` // override global prefs
	Ctype     componentType
}

// initializeInherent attributes of spec and child pathComponents
func (S *spec) initializeInherent() {
	S.Ctype = cfgComponent
	for _, src := range S.Sources {
		src.Parent = S.Alias
		src.Ptype = sourceComponent
	}
	for _, tgt := range S.Targets {
		tgt.Parent = S.Alias
		tgt.Ptype = targetComponent
	}
}

func (S spec) allInitialized() bool {
	all := S.Ctype > 0
	for _, src := range S.Sources {
		all = all && src.Ctype > 0 && src.Parent != ""
	}
	for _, tgt := range S.Sources {
		all = all && tgt.Ctype > 0 && tgt.Parent != ""
	}
	return all
}

func (S spec) getAlias() string        { return S.Alias }
func (S spec) getCtype() componentType { return S.Ctype }

func (S *spec) getSource(alias string) *pathComponent {
	for _, src := range S.Sources {
		if alias == src.Alias {
			return &src
		}
	}
	return nil
}
func (S *spec) getTarget(alias string) *pathComponent {
	for _, tgt := range S.Targets {
		if alias == tgt.Alias {
			return &tgt
		}
	}
	return nil
}

func (S spec) status() string {
	expln := fmt.Sprintf("spec:'%s' - Sources:\n%+v", S.Alias, S.Sources)
	return expln
}

// ╭─────────────────────────────────────────────────────────╮
// │                 MAKING/RUNNING COPYJOBS                 │
// ╰─────────────────────────────────────────────────────────╯
func (S spec) runCopy() error {
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

// NOTE: FROM GLOBALMODIFY.GO ──────────────────────────────────────────────────
func (S *spec) GetIgnores() *[]string { return &S.Ignorepat }
func (S *spec) GetLocalPrefs() *prefs { return &S.Overrides }

// func (S *spec) getChildByPath(path string) *pathComponent { }
// TODO: implement (in globalModify.go)
func (S spec) IsPathChild(path string) bool {
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

func (S spec) GetIfChild(path string) *pathComponent {
	for _, src := range S.Sources {
		if src.Alias == path || src.Path == pops.MakeAbs(path) {
			return &src
		}
	}
	for _, tgt := range S.Targets {
		if tgt.Alias == path || tgt.Path == pops.MakeAbs(path) {
			return &tgt
		}
	}
	return nil
}

func (S *spec) AddIgnores(ignores []string) {
	S.Ignorepat = append(S.Ignorepat, ignores...)
}
func (S *spec) CheckAddPath(path string, isSource bool) bool {
	var ctyp componentType
	if isSource {
		ctyp = sourceComponent
	} else {
		ctyp = targetComponent
	}
	if !S.IsPathChild(path) {
		S.Sources = append(S.Sources, *newPathComponent(path, ctyp))
		return true
	}
	return false
}

// CheckAddMultiplePaths adds paths to spec.Sources if isSource, spec.Targets if !isSource
func (S *spec) CheckAddMultiplePaths(paths []string, isSource bool) {
	var ctyp componentType
	if isSource {
		ctyp = sourceComponent
	} else {
		ctyp = targetComponent
	}
	_ = ctyp
	for _, p := range paths {
		_ = p //if pc := S.getChildByPath(p); pc == nil { } else { }
	}
}
