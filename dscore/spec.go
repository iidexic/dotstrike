package dscore

import (
	"fmt"
	"maps"
	"slices"
	"strconv"

	pops "iidexic.dotstrike/pathops"
	"iidexic.dotstrike/uout"
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
	Sources    []PathComponent `toml:"sources"`    // paths marked as origin points
	Targets    []PathComponent `toml:"targets"`    // paths  marked as destination points
	Ignorepat  preIgnoreList   `toml:"ignores"`    // ignorepat that apply to all sources
	OverrideOn bool            `toml:"overrideOn"` // enable overrides, prevent Overrides being over-written
	Overrides  prefs           `toml:"overrides"`  // override global prefs
	Ctype      componentType   // I thought this got deleted
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

// ── Find/Get Spec Info ──────────────────────────────────────────────
func specEqual(S, S2 Spec) bool {
	return S.Alias == S2.Alias && S.Overrides.equal(S2.Overrides) && S.Ctype == S2.Ctype &&
		slices.EqualFunc(S.Sources, S2.Sources, pathComponentEqual) &&
		slices.EqualFunc(S.Targets, S2.Targets, pathComponentEqual) &&
		slices.Equal(S.Ignorepat, S2.Ignorepat)
}

func (S Spec) Identify() string { return S.Alias }

func (S *Spec) Detail() string {
	out := uout.NewOutf("- SPEC: %s ------", S.Alias)
	out.IndR()
	if len(S.Sources)+len(S.Targets) == 0 {
		out.V("No Path Components (0 Sources, 0 Targets)")
	} else {
		out.F("Sources (%d):", len(S.Sources))
		out.IndR().ILV(S.Sources)
		out.IndL().F("Targets (%d):", len(S.Targets))
		out.IndR().ILV(S.Targets)
		out.IndL()
	}
	out.IfV(S.OverrideOn, "Config Overrides (enabled)", "Config Overrides (disabled)")
	lcfg := len(S.Overrides.Bools)
	out.AF(" (%d options set)", lcfg)
	out.IndR()
	for opt, b := range S.Overrides.Bools {
		out.IfF(b, "opt %s: On", "opt %s: Off", opt, opt)
	}
	out.IndL()
	if igN := len(S.Ignorepat); igN > 0 {
		out.F("%d Ignore Patterns: ", igN)
		out.FlatLV(S.Ignorepat)
	}
	return out.String()
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

func (S Spec) String() string {
	return S.Detail()
}

func (S *Spec) DetailSources(parentName bool) string {
	if len(S.Sources) == 0 {
		return "	none"
	}
	out := uout.NewOut("Sources:")
	out.IndR()
	out.ILV(S.Sources)
	return out.String()
}
func (S *Spec) DetailTargets(parentName bool) string {
	if len(S.Targets) == 0 {
		return "	none"
	}
	out := uout.NewOut("Targets:")
	out.IndR()
	out.ILV(S.Targets)
	return out.String()
}

// func (S *Spec) GetIgnores() *[]string { return &S.Ignorepat }

func (S *Spec) GetLocalPrefs() *prefs { return &S.Overrides }

// SetOverrideMap sets the override map for the spec. returns a slice of failed options
func (S *Spec) SetOverrideMap(mpref map[string]bool) ([]string, error) {
	fails := make([]string, len(mpref))
	n := 0
	var eout error
	for k, b := range mpref { // all this to intercept overrideOn
		if opt := OptionID(k); opt == BoolOverrideOn {
			tempData.Modify()
			S.OverrideOn = b
		} else if opt != NotAnOption {
			//WARN: never errors for now but may need to check in future
			// oh it actually can error although still don't think it is able as set
			e := S.Overrides.setOpt(opt, b)
			if eout == nil {
				eout = e
			} else {
				eout = fmt.Errorf("%w\n%w", eout, e)
			}
		} else {
			fails[n] = k
			n++
		}

	}
	return fails[:n], eout
}

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

func (S *Spec) GetExistingChildren(identifiers []string) []*PathComponent {
	components := make([]*PathComponent, 0, len(identifiers))
	for _, id := range identifiers {
		if pc := S.GetIfChild(id); pc != nil {
			components = append(components, pc)
		}
	}
	return components
}

// Does this work right now??
func (S *Spec) GetMatchingComponents(identifiers []string, isSource bool) []*PathComponent {
	lenids := len(identifiers)
	var cmpExisting []PathComponent
	cmpMatching := make([]*PathComponent, lenids)
	if isSource {
		cmpExisting = S.Sources
	} else {
		cmpExisting = S.Targets
	}
	sm := make(map[string]struct{}, lenids)
	for _, sid := range identifiers {
		// This shit is a mess
		key := pops.MakeAbs(sid)
		sm[key] = struct{}{}
	}
	n := 0
	for i, comp := range cmpExisting {
		if _, ok := sm[comp.Path]; ok {
			cmpMatching[n] = &cmpExisting[i]
			n++
		} else if comp.Alias != "" {
			if _, ok := sm[comp.Alias]; ok {
				cmpMatching[n] = &cmpExisting[i]
				n++
			}
		}
	}
	if n < lenids {
		cmpMatching = cmpMatching[:n]
	}
	return cmpMatching
}

// GetIfChild returns a pointer to the child source or target with the path or alias passed. Returns nil if none found
func (S *Spec) GetIfChild(identifier string) *PathComponent {
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

// TODO: Replace CheckAddPath with S.AddSource! why
func (S *Spec) AddSource(path string, ignorelist ...string) error {
	if !S.IsPathChild(path) {
		tempData.Modify()
		abs := pops.MakeAbs(path)
		S.Sources = append(S.Sources,
			PathComponent{
				Path:    path,
				Abspath: abs,
				Ctype:   sourceComponent,
				Parent:  S.Alias,
			})
		return nil
	}
	echild := S.GetIfChild(path)
	return fmt.Errorf("path `%s` is already in spec %s as a %s!", path, S.Alias, echild.Ctype.String())
}
func (S *Spec) AddIgnores(ignores []string) {
	S.Ignorepat = append(S.Ignorepat, ignores...)
}

func (S *Spec) CheckAddPath(path string, isSource bool) bool {
	if !S.IsPathChild(path) {
		path = pops.TildeExpand(path)
		tempData.Modify()
		if isSource {
			S.Sources = append(S.Sources, *newPathComponent(path, sourceComponent))
		} else {
			S.Targets = append(S.Targets, *newPathComponent(path, targetComponent))
		}
		return true
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

func (S *Spec) DeleteByPtr(components ...*PathComponent) error {
	mc := make(map[string]bool, len(components))
	for _, c := range components {
		mc[c.Path] = false
	}
	for i, src := range S.Sources {
		if _, ok := mc[src.Path]; ok {
			S.removeSourceByIndex(i)
		}
	}

	for i, tgt := range S.Targets {
		if _, ok := mc[tgt.Path]; ok {
			S.removeTargetByIndex(i)
		}
	}
	return nil
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
	var oplist []PathComponent
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

func (S *Spec) getComponentList(isSource bool) *[]PathComponent {
	if isSource {
		return &S.Sources
	}
	return &S.Targets
}

// WipeComponentList deletes everything from Sources if isSource, or Targets if !isSource
func (S *Spec) WipeComponentList(isSource bool) {
	tempData.Modify()
	if isSource {
		S.Sources = make([]PathComponent, 0)
	} else {

		S.Targets = make([]PathComponent, 0)
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
