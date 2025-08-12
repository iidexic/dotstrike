package dscore

import (
	"fmt"
	"slices"
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
// Methods:
//   - Detail(): nice printable spec detail string
//   - GetLocalPrefs(): returns Overrides (dscore.prefs)
//   - GetIgnores(): returns []string of ignore patterns
//   - IsPathSource/IsPathTarget(path string): checks for existing child component with path
//   - GetIfChild: If path is within a child component, a pointer to that component is returned
type Spec struct {
	Alias      string          `toml:"alias"`      // name, unique
	Sources    []pathComponent `toml:"sources"`    // paths marked as origin points
	Targets    []pathComponent `toml:"targets"`    // paths  marked as destination points
	Ignorepat  []string        `toml:"ignores"`    // ignorepat that apply to all sources
	OverrideOn bool            `toml:"overrideOn"` // enable overrides, prevent Overrides being over-written
	Overrides  prefs           `toml:"overrides"`  // override global prefs
	Ctype      componentType
}

var ErrID error = fmt.Errorf("Identifier not found")
var ErrComponentType = fmt.Errorf("Identifier found with wrong component type")

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
func (S *Spec) allInitialized() bool {
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

func (S *Spec) getAlias() string        { return S.Alias }
func (S *Spec) getCtype() componentType { return S.Ctype }

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
	lines = append(lines, "Spec: "+S.Alias, "-------------------",
		"Sources:", S.DetailSources(false), "Targets:", S.DetailTargets(false))

	// ── Sources ──────────────────────────────
	//lines = append(lines, "Sources:")
	// if len(S.Sources) == 0 {
	// 	lines = append(lines, "	none")
	// }
	// for i, src := range S.Sources {
	// 	lines = append(lines, fmt.Sprintf("[src %d]", i+1))
	// 	lines = append(lines, src.Detail())
	// }

	// ── Targets ──────────────────────────────
	// lines = append(lines, "Targets:")
	// if len(S.Targets) == 0 {
	// 	lines = append(lines, "	none")
	// }
	// for i, tgt := range S.Targets {
	// 	lines = append(lines, fmt.Sprintf("[tgt %d]", i+1))
	// 	lines = append(lines, tgt.Detail())
	// }

	// ── Overrides ────────────────────────────
	lines = append(lines, fmt.Sprintf("Overrides Enabled: %t", S.OverrideOn))
	if !S.Overrides.equal(gd.data.Prefs) {
		lines = append(lines, fmt.Sprintf(`Overrides:
	Keep Repo: %t
	Keep Hidden Files: %t
	Use Global Target: %t`, S.Overrides.KeepRepo, S.Overrides.KeepHidden, S.Overrides.GlobalTarget))
	}

	// ── Ignores ──────────────────────────────
	if len(S.Ignorepat) > 0 {
		lines = append(lines, "Ignore Patterns:")
		for i, pat := range S.Ignorepat {
			lines = append(lines, fmt.Sprintf("	 %d.) '%s'", i, pat))
		}
	}
	lines = append(lines, "")
	return strings.Join(lines, "\n")
}
func (S *Spec) ShortDetail() string {
	line := fmt.Sprintf("%s, ", S.Alias)
	sl := len(S.Sources)
	tl := len(S.Targets)
	switch sl {
	case 1:
		line += fmt.Sprintf("[1 src: %s]", S.Sources[0].Path)
	default:
		line += fmt.Sprintf("[%d sources]", sl)
	}
	switch tl {
	case 1:
		line += fmt.Sprintf("[1 tgt: %s]", S.Targets[0].Path)
	default:
		line += fmt.Sprintf("[%d targets]", tl)
	}

	if S.OverrideOn {
		line += "(overrides on)"
	}
	return line
}

func (S *Spec) DetailSources(parentName bool) string {
	if len(S.Sources) == 0 {
		return "	none"
	}
	ss := make([]string, len(S.Sources))
	for i, src := range S.Sources {
		sstr := ""
		if parentName {
			sstr = fmt.Sprintf("spec %s | ", S.Alias)
		} else {
			sstr = "	"
		}
		if src.Alias != "" {
			sstr += fmt.Sprintf("%s: ", src.Alias)
		} else {
			sstr += fmt.Sprintf("[%d]:", i)
		}
		sstr += src.Path
		if src.Abspath != src.Path && src.Abspath != "" {
			sstr += fmt.Sprintf(" (%s)", src.Abspath)
		} else if src.Abspath == "" {
			sstr += "(WARNING: NO ABSOLUTE PATH)"
		}
		if len(src.Ignores) > 0 {
			sstr += fmt.Sprintf("\n		ignores:%v", src.Ignores)
		}
		ss[i] = sstr
	}
	return strings.Join(ss, "\n")
}
func (S *Spec) DetailTargets(parentName bool) string {
	if len(S.Targets) == 0 {
		return "	none"
	}
	ss := make([]string, len(S.Targets))
	for i, tgt := range S.Targets {
		sstr := ""
		if parentName {
			sstr = fmt.Sprintf("spec %s | ", S.Alias)
		} else {
			sstr = "	"
		}
		if tgt.Alias != "" {
			sstr += fmt.Sprintf("%s: ", tgt.Alias)
		} else {
			sstr += fmt.Sprintf("[%d]:", i)
		}
		sstr += tgt.Path
		if tgt.Abspath != tgt.Path && tgt.Abspath != "" {
			sstr += fmt.Sprintf(" (%s)", tgt.Abspath)
		} else if tgt.Abspath == "" {
			sstr += "(WARNING: NO ABSOLUTE PATH)"
		}
		if len(tgt.Ignores) > 0 {
			sstr += fmt.Sprintf("\n		ignores:%v", tgt.Ignores)
		}
		ss[i] = sstr
	}
	return strings.Join(ss, "\n")
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

// IsPathChild looks for the path within the Spec's pathComponent slices
func (S *Spec) IsPathChild(path string) bool {
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

func (S *Spec) IsPathSource(path string) bool {
	for _, src := range S.Sources {
		if src.Alias == path || src.Path == pops.MakeAbs(path) {
			return true
		}
	}
	return false
}

func (S *Spec) IsPathTarget(path string) bool {
	for _, tgt := range S.Targets {
		if tgt.Alias == path || tgt.Path == pops.MakeAbs(path) {
			return true
		}
	}
	return false
}

func (S *Spec) GetExistingChildren(identifiers []string) []*pathComponent {
	components := make([]*pathComponent, 0, len(identifiers))
	for _, id := range identifiers {
		if pc := S.GetIfChild(id); pc != nil {
			components = append(components, pc)
		}
	}
	return components
}

// GetIfChild returns a pointer to the child source or target with the path or alias passed. Returns nil if none found
func (S *Spec) GetIfChild(identifier string) *pathComponent {
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

func (S *Spec) removeSourceByIndex(index int) {
	tempData.Modify()

	if index < len(S.Sources) {
		S.Sources = slices.Delete(S.Sources, index, index+1)
	}
}
func (S *Spec) removeTargetByIndex(index int) {
	tempData.Modify()
	if index < len(S.Targets) {
		S.Targets = slices.Delete(S.Targets, index, index+1)
	}
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
func (S *Spec) AliasIfChild(alias string, identifier string, isSource bool) bool {
	tempData.Modify()
	if pc := S.GetIfChild(identifier); pc != nil {
		if (pc.Ctype == sourceComponent && isSource) ||
			(pc.Ctype == targetComponent && !isSource) {
			pc.setAlias(alias)
		} else {
		}
	}
	return false
}

func (S *Spec) DeleteIfChild(identifier string, isSource bool, singleDelete bool) int {
	tempData.Modify()
	count := 0
	if isSource {
		for i := range S.Sources {
			if S.Sources[i].MatchesID(identifier) {
				S.removeSourceByIndex(i)
				count++
				if singleDelete {
					return count
				}
			} else {
			}
		}
	} else {
		for i := range S.Targets {
			if S.Targets[i].MatchesID(identifier) {
				S.removeTargetByIndex(i)
				count++
				if singleDelete {
					return count
				}
			}
		}

	}
	return count
}

// WipeComponentList deletes everything from Sources if isSource, or Targets if !isSource
func (S *Spec) WipeComponentList(isSource bool) {
	tempData.Modify()
	if isSource {
		S.Sources = make([]pathComponent, 0)
	} else {

		S.Targets = make([]pathComponent, 0)
	}
}

// CheckAddMultiplePaths adds paths to spec.Sources if isSource, spec.Targets if !isSource
// returns slice of bools indicating which indices in paths were added succesfully
func (S *Spec) CheckAddMultiplePaths(paths []string, isSource bool) []bool {
	tempData.Modify()
	b := make([]bool, len(paths))
	for i, p := range paths {
		b[i] = S.CheckAddPath(p, isSource)
	}
	return b
}

// ── Running spec copy jobs ──────────────────────────────────────────
func (S *Spec) RunCopy(global bool) error {
	if !S.allInitialized() {
		return fmt.Errorf("spec not initialized: %s", S.Alias)
	}
	copymachine := pops.GetCopierMaschine()
	//NOTE:
	for y, tgt := range S.Targets {
		for x, src := range S.Sources {
			job := copymachine.NewJob(S.Alias+"."+fmt.Sprintf("%d", x)+"."+fmt.Sprintf("%d", y), src.Path, tgt.Path)
			job.JobOptionMakeSubdir(true)
			_ = job
		}
	}
	return nil
}
