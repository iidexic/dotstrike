package dscore

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	pops "iidexic.dotstrike/pathops"
)

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
		lines = append(lines, fmt.Sprintf("Overrides:\n%s", S.Overrides.Detail()))
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

// Loops through component list and deletes all matching given ID
func (S *Spec) runtimeRemoveMatching(id string, isSource bool) {
	if isSource {
		for i, t := range S.Sources {
			if t.MatchesID(id) {
				S.Sources = slices.Delete(S.Sources, i, i+1)
			}
		}
	} else {
		for i, t := range S.Targets {
			if t.MatchesID(id) {
				S.Targets = slices.Delete(S.Targets, i, i+1)
			}
		}
	}
}
func (S *Spec) removeTargetByIndex(index int) {
	tempData.Modify()
	if index < len(S.Targets) {
		S.Targets = slices.Delete(S.Targets, index, index+1)
	}
}

// ── Modifying Spec Data ─────────────────────────────────────────────

// TODO: Replace CheckAddPath with S.AddSource!
func (S *Spec) AddSource(path string, ignorelist ...string) error {
	if !S.IsPathChild(path) {
		tempData.Modify()
		abs := pops.MakeAbs(path)
		S.Sources = append(S.Sources,
			pathComponent{
				Path:    path,
				Abspath: abs,
				Ctype:   sourceComponent,
				Parent:  S.Alias,
			})
		return nil
	}
	echild := S.GetIfChild(path)
	return fmt.Errorf("path `%s` is already in spec %s as a %s!", path, S.Alias, echild.getCtype().String())
}
func (S *Spec) AddIgnores(ignores []string) {
	S.Ignorepat = append(S.Ignorepat, ignores...)
}

// WARNING: Non-Persistent
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
func (S *Spec) GetMatching(ids []string, isSource bool) []int {
	imap := make(map[int]bool, len(ids))
	for _, id := range ids {
		if i := S.textMatchesComponent(id, isSource); i >= 0 {
			imap[i] = false
		}
	}
	indices := make([]int, len(imap))
	n := 0
	for k := range imap {
		indices[n] = k
		n++
	}
	return indices
}

// TextMatchesComponent identifies first source/target where s matches and returns its index
// The matches tried are:
//   - sv as an integer index in the slice (ie S.Sources[int(sv)])
//   - sv match to source/target alias (ie clean(sv) == source.Alias)
//   - sv match full path or base name (ie clean(sv) == filepath.Base(source.Path))
func (S *Spec) textMatchesComponent(xid string, isSource bool) int {
	var oplist []pathComponent
	if isSource {
		oplist = S.Sources
	} else {
		oplist = S.Targets
	}

	index, e := strconv.Atoi(xid)
	if e == nil && index < len(oplist) {
		return index
	}
	// now for the other ones
	xid = QuickClean(xid)
	for i := range oplist {
		if oplist[i].MatchesID(xid) {
			return i
		}
	}
	return -1
}

func (s *Spec) cloneSelf() *Spec {
	new := &Spec{Alias: s.Alias, OverrideOn: s.OverrideOn, Ctype: specComponent}
	copy(new.Sources, s.Sources)
	copy(new.Targets, s.Targets)
	if len(s.Ignorepat) > 0 {
		copy(new.Ignorepat, s.Ignorepat)
	}
	if len(s.Overrides.Bools) > 0 {
		maps.Copy(new.Overrides.Bools, s.Overrides.Bools)
	}
	return new
}

func (S *Spec) getComponentList(isSource bool) *[]pathComponent {
	if isSource {
		return &S.Sources
	}
	return &S.Targets
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

func (s *Spec) sourcePaths() []string {
	paths := make([]string, len(s.Sources))
	for i, src := range s.Sources {
		if src.Abspath != "" {
			paths[i] = src.Abspath
		} else {
			paths[i] = src.Path
		}
	}
	return paths
}

func (s *Spec) targetPaths() []string {
	paths := make([]string, len(s.Targets))
	for i, tgt := range s.Targets {
		if tgt.Abspath != "" {
			paths[i] = tgt.Abspath
		} else {
			paths[i] = tgt.Path
		}
	}
	return paths
}

func (S *Spec) stripComponentList(ikeep []int, isSource bool) {
	oplist := S.getComponentList(isSource)
	newlist := KeepIndices(*oplist, ikeep)
	oplist = &newlist
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

// ── Running spec copy jobs(MOVE TO EXECUTE) ──────────────────────────────────────────
// func (S *Spec) RunCopy(global bool) error {
// 	if !S.allInitialized() {
// 		return fmt.Errorf("spec not initialized: %s", S.Alias)
// 	}
// 	copymachine := pops.GetCopierMaschine()

// 	for y, tgt := range S.Targets {
// 		for x, src := range S.Sources {
// 			job := copymachine.NewJob(S.Alias+"."+fmt.Sprintf("%d", x)+"."+fmt.Sprintf("%d", y), src.Path, tgt.Path)
// 			job.JobOptionMakeSubdir(true)
// 			_ = job
// 		}
// 	}
// 	return nil
// }
