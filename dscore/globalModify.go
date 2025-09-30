package dscore

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	pops "iidexic.dotstrike/pathops"
)

var ErrNotUnique error = errors.New("Attempted to set component alias to a non-unique value")
var ErrParentNotFound error = errors.New("Component.Parent did not match any existing alias.")
var ErrAliasNotFound error = errors.New("Component.Parent did not match any existing alias.")
var ErrBadKey error = errors.New("Provided map key does not exist")

type Temp interface {
	NewSpec(string, ...string) (*Spec, error)
	DeleteSpec(*Spec) bool
	Select(string) bool
	SelectedSpec() *Spec
	GetSpec(string) *Spec
}

// TODO: add to cmd-spec in place of NewSpec
func (gm *globalModify) NewSpecEmpty(alias string) (*Spec, error) {
	cAlias, err := uniqueSpecAlias(alias)
	if err != nil {
		return nil, err
	}
	gm.Modify()
	gm.Specs = append(gm.Specs, Spec{Alias: cAlias, Ctype: specComponent})
	return &gm.Specs[len(gm.Specs)-1], nil
}

// If user requests a new spec but provides an existing alias,
// uniqueSpecAlias will add an incrementing integer to ensure new spec alias is unique
// limit hard set at 10 attempts before erroring
func uniqueSpecAlias(alias string) (string, error) {
	alias = standardizeAlias(alias)
	if tempData.GetSpec(alias) != nil {
		n := 1
		for {
			ainc := fmt.Sprintf("%s%d", alias, n)
			if tempData.GetSpec(ainc) == nil {
				return ainc, nil
			}
			if n > 9 {
				return "", fmt.Errorf("spec '%s' + incrementals already exist", alias)
			}
		}
	}
	return alias, nil
}

// NewSpec makes a new spec. It is automatically added to TempData.
// Adds paths in src and tgt as new sources and targets, respectively.
// Error occurs if the provided alias is already in use and increment limit is reached
func (gm *globalModify) NewSpec(alias string, src, tgt []string) (*Spec, error) {
	cAlias, err := uniqueSpecAlias(alias)
	if err != nil {
		return nil, fmt.Errorf("Alias already in use (%s, standardized = %s)",
			alias, standardizeAlias(alias))
	}
	s := Spec{Alias: cAlias, Ctype: specComponent}
	if !gm.initialized {
		return nil, ErrNoInit
	}
	if len(src) > 0 {
		s.CheckAddMultiplePaths(src, true)
	}
	if len(tgt) > 0 {
		s.CheckAddMultiplePaths(tgt, false)
	}
	gm.Specs = append(gm.Specs, s)
	newSpecPtr := &gm.Specs[len(gm.Specs)-1] //works
	gm.Modified = true
	return newSpecPtr, nil
}

// adds paths as components without causing data to be written to user file
//
// temp components will still be written if another change is made
// that triggers tempData.Modified() within the same run.
// This includes all public functions/methods that make user data changes.
// Currently there is no single run that can make both temp/persistent changes.
func (s *Spec) temporaryComponents(isSource bool, paths ...string) error {
	etext := make([]string, 0, len(paths))
	if isSource {
		for _, p := range paths {
			if !s.IsPathChild(p) {
				s.Sources = append(s.Sources, *newPathComponent(p, sourceComponent))
			} else {
				etext = append(etext, p)
			}
		}
	} else {
		for _, p := range paths {
			if !s.IsPathChild(p) {
				s.Sources = append(s.Sources, *newPathComponent(p, sourceComponent))
			} else {
				etext = append(etext, p)
			}
		}
	}
	if numer := len(etext); numer > 0 {
		return fmt.Errorf("Failed to add %d paths\n(%s)", numer, strings.Join(etext, ", "))
	}
	return nil
}

// GetSpec searches the spec list and returns the *spec that matches provided alias
//
// If no matches are found, returns nil.
//
// The lookup is exact to the passed alias string
func (gm *globalModify) GetSpec(alias string) *Spec {
	alias = standardizeAlias(alias)
	for i := range gm.Specs {
		if gm.Specs[i].Alias == alias {
			return &gm.Specs[i]
		}
	}
	return nil
}

// GetSpecs finds and returns all Spec values whose unique alias is in aliases.
// It also returns all indices in aliases that found no match.
func (gm globalModify) GetSpecs(forceSelected bool, aliases ...string) ([]*Spec, []int) {
	aqty := len(aliases)
	// TODO:(low)  eliminate the nested loop
	if aqty == 0 {
		//WARNING: Probably bad idea to return nil slice instead of an empty one, not sure
		return []*Spec{gm.SelectedSpec()}, nil
	}
	sel, selIn := gm.SelectedSpec(), false
	specs := make([]*Spec, aqty+1)
	notfound := make([]int, 0, aqty)
	n := 0
	for i, a := range aliases {
		if specs[n] = gm.GetSpec(a); specs[n] != nil {
			n++
			if selIn || specs[n] == sel {
				selIn = true
			}
		} else {
			notfound = append(notfound, i)
		}
	}
	if forceSelected && !selIn {
		specs[n] = sel
		n++
	}
	if n < aqty+1 {
		specs = specs[:n]
	}
	if len(notfound) == 0 {
		return specs, nil
	}
	return specs, notfound
}

// DeleteSpec deletes the spec *sptr (persistent).
func (gm *globalModify) DeleteSpec(sptr *Spec) bool {
	for i := range gm.Specs {
		if spec := &gm.Specs[i]; spec == sptr {
			gm.Modify()
			if isLastAndSelectedSpec(i) {
				ResetSpecSelection()
			}
			// Does this cause a problem if given
			gm.Specs = slices.Delete(gm.Specs, i, i+1)
			return true
		}
	}
	return false
}

func isLastAndSelectedSpec(i int) bool { return i+1 == len(tempData.Specs) && i == tempData.Selected }

// ResetSpecSelection Resets the selected spec to 0 (persistent).
func ResetSpecSelection() { tempData.Modify(); tempData.Selected = 0 }

// Modify toggles gm modified bool to true; this is neccessary to write userdata changes to file.
// Modify will generally be placed directly before the code chunk that actually makes the change.
func (gm *globalModify) Modify() { gm.Modified = true }

// TODO: Update SetOptionBool to handle nonexistant values.
//	Check what using on?? Would not be simple.

// TODO: (Hi-Mid fix) ensure that global prefs gets populated with falses where values not entered by user.

// Sets the global option opt. Persistent
func (gm *globalModify) SetOptionBool(opt ConfigOption, newValue bool) bool {
	val, exist := gm.Prefs.Bools[opt]
	switch {
	case !exist:
		fallthrough
	case exist && val != newValue:
		gm.Modify()
		gm.Prefs.Bools[opt] = newValue
		return true
	case exist:
		return true
	}
	return false
}

// SetOptionString sets global prefs[opt] to newValue
// Returns an error if opt is not a string option
func (gm *globalModify) SetOptionString(opt ConfigOption, newValue string) error {
	if !opt.IsString() {
		return fmt.Errorf("Not a string option,")
	}
	switch opt {
	case StringGlobalTargetPath:
		tempData.Modify()
		newpath, e := pops.Abs(newValue)
		if e != nil {
			return e
		}
		gm.GlobalTargetPath = newpath
		return nil
	}
	return ErrID
}

func (gm *globalModify) Select(alias string) bool {

	index := gm.globalData.findAliasIndex(alias)
	if index < 0 {
		return false
	}
	if index != gm.globalData.Selected {
		gm.Modified = true
		gm.globalData.Selected = index
	}
	return true
}
func (gm *globalModify) SelectPtr(spec *Spec) bool {
	for i := range gm.Specs {
		if &gm.Specs[i] == spec {
			gm.Modified = true
			gm.globalData.Selected = i
			return true
		}
	}
	return false
}

func (gm *globalModify) ChangeSpecAlias(spec *Spec, newAlias string) bool {
	newAlias = standardizeAlias(newAlias)
	otherspec := gm.GetSpec(newAlias)
	if otherspec != nil {
		return false
	}
	gm.Modify()
	spec.Alias = newAlias
	return true
}

func (gm *globalModify) specByIndex(i int) *Spec {
	if i < len(gm.Specs) {
		return &gm.globalData.Specs[i]
	}
	return nil
}

// TODO:(mid-refactor/system) re-work SetSpecOverrides/setOptMap to take into account OverrideOn and other non-prefs settings.
// Should function similarly or same to setting GlobalTargetPath

func (gm *globalModify) SetSpecOverridesMap(s *Spec, newValues map[string]bool) []string {
	fails, e := s.Overrides.setOptMap(newValues)
	if e != nil {
		for i, f := range fails {
			// Check to see if
			if matchOverrideOn(f) {
				s.OverrideOn = newValues[f]
				fails = slices.Delete(fails, i, i+1)
				break
			}
		}
	}
	return fails
}

func matchOverrideOn(t string) bool { return strings.Contains("override", strings.ToLower(t)) }

func (gm *globalModify) SetSpecEnableOverrides(s *Spec, enable bool) bool {
	if s.OverrideOn != enable {
		s.OverrideOn = enable
		return true
	}
	return false
}
func (gm *globalModify) SelectedSpec() *Spec { return gm.specByIndex(gm.Selected) }

// findAliasIndex searches for alias within []Specs
// probably remove
func (gd *globalData) findAliasIndex(alias string) int {
	for i := range gd.Specs {
		if alias == gd.Specs[i].Alias {
			return i
		}
	}
	return -1
}

func (gm *globalModify) CountComponents() int {
	count := len(gm.Specs)
	for i := range gm.Specs {
		count += len(gm.Specs[i].Sources) + len(gm.Specs[i].Targets)
	}
	return count
}

// setOptMap will modify all prefs/overrides with a key assigned in mpref.
// keys are not case-sensitive, and all spaces are removed.
// Accepted keys are in dsconfig.go or config package
//
// Returns list of strings that failed to correspond to an option
func (p *prefs) setOptMap(mpref map[string]bool) ([]string, error) {
	fails := make([]string, 0, len(mpref))
	var ferr error
	for k, b := range mpref {
		err := p.setByName(k, b)
		if err != nil {
			fails = append(fails, k)
			if ferr == nil {
				ferr = fmt.Errorf("%w", err)
			} else {
				ferr = fmt.Errorf("%w\n%w", ferr, err)
			}
		} else {
			//IDEA: Try stripping map of values that have been written with no return.
		}
	}
	return fails, ferr
}

// TODO:(low-refactor) clean up the SetOption mess.
func (p *prefs) setByName(name string, val bool) error {
	if opt := OptionID(name); opt != NotAnOption {
		e := p.setOpt(opt, val)
		if e != nil {
			return fmt.Errorf("prefs.setOpt error 0 name='%s',OptionID='%s';\n%w", name, opt.String(), e)
		}
		return nil
	}
	return fmt.Errorf("OptionID: String %s produced NotAnOption", name)
}

func (p *prefs) setOpt(opt ConfigOption, val bool) error {
	if opt.IsBool() && opt.IsRealOption() {
		p.Bools[opt] = val
		tempData.Modify()
		return nil
	}
	failString := "error assigning option" + opt.String() + " ("
	if !opt.IsBool() {
		failString += "not a boolean option, "
	}
	if !opt.IsRealOption() {
		failString += "not real option"
	}
	failString += ")"
	return fmt.Errorf("%s", failString)
}

func (gd *globalData) SetGlobalTargetPath(path string) { gd.GlobalTargetPath = pops.CleanPath(path) }
