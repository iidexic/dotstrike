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

//TODO: standardize modify. Now prefer to put as deep/as close to the actual modify as possible,
// Intent to avoid write to file as much as possible

// TODO: Replace NewSpec.  The weird sources-target thing with paths is bad.
// Also, there are flags for that...
func (gm *globalModify) NewSpec(alias string, paths ...string) (*Spec, error) {
	s := Spec{Alias: alias, Ctype: specComponent}
	if !gm.initialized {
		return nil, ErrNoInit
	}
	switch {
	case len(paths) > 1:
		s.CheckAddPath(paths[len(paths)-1], false)
		s.addSources(paths[:len(paths)-1]...)
	case len(paths) == 1:
		s.CheckAddPath(paths[0], true)
	}
	gm.Specs = append(gm.Specs, s)
	newSpecPtr := &gm.Specs[len(gm.Specs)-1] //works
	gm.Modified = true
	return newSpecPtr, nil
}
func (s *Spec) addSources(paths ...string) []bool {
	added := make([]bool, len(paths))
	for i, src := range paths {
		added[i] = s.CheckAddPath(src, true)
	}
	return added
}

// adds paths as components without causing data to be written to user file
//
// NOTE: if another change is made that triggers tempData.Modified(), these will still be written if not manually corrected.
// This includes all public functions/methods that make user data changes.
// IF NEED THIS TO NOT HAPPEN: Need a new component struct OR break out path components into types
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

// GetModifiableSpec returns a pointer to a spec that will be encoded when the program exits
// If spec exists but has not yet been modified, this adds that spec to gm.Specs
// TODO: Update/Replace this as it is now entirely unnecessary. Also, the errors are a mess
func (gm *globalModify) GetModifiableSpec(alias string) (*Spec, error) {
	ermsg := make([]string, 1, len(gm.Specs)+len(gd.data.Specs))
	ermsg[0] = "[MODIFY_SPECS]"
	for i, s := range gm.Specs {
		if s.Alias == alias {
			gm.Modified = true
			return &gm.Specs[i], nil
		} else {
			ermsg = append(ermsg, s.Alias)
		}
	}
	ermsg = append(ermsg, "[GLOBALDATA_SPECS]")
	for _, s := range gd.data.Specs {
		if s.Alias == alias {
			gm.Specs = append(gm.Specs, s)
			gm.Modified = true
			return &gm.Specs[len(gm.Specs)-1], nil
		} else {
			ermsg = append(ermsg, s.Alias)
		}
	}

	return nil, fmt.Errorf("No matching alias found in:\n%s", strings.Join(ermsg, "\n"))
}

// GetSpec searches the spec list and returns the *spec that matches provided alias
//
// If no matches are found, returns nil.
//
// The lookup is exact to the passed alias string
func (gm *globalModify) GetSpec(alias string) *Spec {
	for i := range gm.Specs {
		if gm.Specs[i].Alias == alias {
			return &gm.Specs[i]
		}
	}
	return nil
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

func (gm *globalModify) Modify() { gm.Modified = true }

/* Unused
// SetNamedOptionBool performs option lookup and sets selected configOption opt to newValue (persistent).
// Returns true if a config value was changed, false otherwise
 func (gm *globalModify) SetNamedOptionBool(optName string, newValue bool) bool {
	//BUG panics on optName->NotAnOption
	return gm.SetOptionBool(OptionID(strings.ToLower(strings.TrimSpace(optName))), newValue)
}

func (gm *globalModify) SetNamedOptionString(optName string, newValue string) error {
	return gm.SetOptionString(OptionID(strings.ToLower(strings.TrimSpace(optName))), newValue)
}
*/

func (gm *globalModify) SetOptionBool(opt ConfigOption, newValue bool) bool {
	val, exist := gm.Prefs.Bools[opt]
	switch {
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

// TODO:(V0.0.1) Delete
// func (p *prefs) OverwriteRaw(newp prefs) error {
// 	if !p.equal(newp) {
// 		for co := range p.bools {
// 			delete(p.bools, co)
// 		}
// 		maps.Copy(p.bools, newp.bools)
// 		// for k, v := range newp.bools {
// 		// 	p.bools[k] = v
// 		// }
//
// 	}
// 	return nil
// }

// setAlias sets the PathComponent alias.
// If PathComponent is not unique, alias is not set, and ErrNotUnique is returned
func (pc *pathComponent) setAlias(alias string) error {
	cfptr := gd.data.getSpec(pc.Parent)
	if cfptr == nil {
		return ErrParentNotFound
	}
	var existingpc *pathComponent
	if pc.Ctype == sourceComponent {
		existingpc = cfptr.getSource(alias)
	} else {
		existingpc = cfptr.getTarget(alias)
	}
	if existingpc != nil {
		return ErrNotUnique
	}
	pc.Alias = alias
	return nil
}
