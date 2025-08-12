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

// for reference:
// type component interface
// 	getAlias() string
// 	getCtype() componentType
// type pathComponent struct
// 	   string
// 	Alias, Abspath, Path, Parent -> string
// 	Ignores -> []string,	Ptype -> pathType,	Ctype -> componentType
// type cfg struct
// 	Alias string, Ignorepat []string
// 	Sources, Targets  []pathComponent
// 	Overrides prefs, Ctype componentType

// NewSpec creates and adds a spec to globalModify/temp data.
//
// If paths are passed, they are added to the spec as sources and targets,
// depending on the length of the paths slice and the index of each path:
//
//   - If paths has exactly 1 item, it is added as a source to the new spec.
//   - If paths has multiple items, the last path is added as a target,
//     and every other item in paths is added as a source.
//
// # Example:
//
//	spec := NewSpec("documents", "C:\foo1", "C:\foo2", "C:\foo3")
//	len(spec.Sources) == 2 , len(spec.Targets) == 1
//	spec.Sources[0].Path == "C:\foo1"
//	spec.Sources[1].Path == "C:\foo2"
//	spec.Targets[0].Path == "C:\foo3"

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
func (gm *globalModify) DeleteSpec(sptr *Spec) bool {
	for i := range gm.Specs {
		if &gm.Specs[i] == sptr {
			// Does this cause a problem if given
			gm.Specs = slices.Delete(gm.Specs, i, i+1)
			return true
		}
	}
	return false
}

type ConfigOption int

const (
	_ ConfigOption = iota
	OptBoolUseGlobalTarget
	OptBoolKeepRepo
	OptBoolKeepHidden
	OptStringGlobalTargetPath
)

func (c ConfigOption) Text() string {
	switch c {
	case OptBoolKeepHidden:
		return "KeepHidden"
	case OptBoolKeepRepo:
		return "KeepRepo"
	case OptBoolUseGlobalTarget:
		return "UseGlobalTarget"
	case OptStringGlobalTargetPath:
		return "SetGlobalTargetPath"
	}
	return ("NotAnOption")
}

func OptionID(optName string) ConfigOption {
	switch {
	case slices.Contains(PrefNameKeepRepo, optName):
		return OptBoolKeepRepo
	case slices.Contains(PrefNameKeepHidden, optName):
		return OptBoolKeepHidden
	case slices.Contains(PrefNameUseGlobalTarget, optName):
		return OptBoolUseGlobalTarget
	case slices.Contains(PrefNameGlobalTargetPath, optName):
		return OptStringGlobalTargetPath

	}
	return 0
}

func (gm *globalModify) Modify() { gm.Modified = true }

// SetOptionBool sets selected configOption opt to newValue.
// Returns true if a config value was changed, false otherwise
func (gm *globalModify) SetNamedOptionBool(optName string, newValue bool) bool {
	return gm.SetOptionBool(OptionID(strings.ToLower(strings.TrimSpace(optName))), newValue)
}

func (gm *globalModify) SetNamedOptionString(optName string, newValue string) error {
	return gm.SetOptionString(OptionID(strings.ToLower(strings.TrimSpace(optName))), newValue)
}

func (gm *globalModify) SetOptionBool(opt ConfigOption, newValue bool) bool {
	switch opt {
	case OptBoolUseGlobalTarget:
		tempData.Modify()
		gm.Prefs.GlobalTarget = newValue
		return true
	case OptBoolKeepRepo:
		tempData.Modify()
		gm.Prefs.KeepRepo = newValue
		return true
	case OptBoolKeepHidden:
		tempData.Modify()
		gm.Prefs.KeepHidden = newValue
		return true
	}

	return false
}

// SetOptionString
func (gm *globalModify) SetOptionString(opt ConfigOption, newValue string) error {
	switch opt {
	case OptStringGlobalTargetPath:
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

// SetM will modify all prefs/overrides with a key assigned in mpref.
// keys are not case-sensitive, and all spaces are removed.
// Accepted keys are:
//   - Keep Git Repo: keeprepo | keep-repo | keep_repo | repo
//   - Keep other hidden: keephidden | keep-hidden | keep_hidden |  hidden
//   - Use Global Target:  globaltarget | global-target | global_target | globaltgt
//
// Returns ErrBadKey if the key does not match an acceptable input. Otherwise, returns nil
func (p *prefs) SetM(mpref map[string]bool) error {
	var fails string
	for k, b := range mpref {
		e := p.SetByName(k, b)
		if e != nil {
			fails = fails + ", " + k
		}
	}
	if len(fails) > 0 {
		return fmt.Errorf("failures: (%s) - error %w", fails, ErrBadKey)
	}
	return nil
}

var PrefNameKeepRepo = []string{"keeprepo", "keep-repo", "keep_repo", "repo"}
var PrefNameKeepHidden = []string{"keephidden", "keep-hidden", "keep_hidden", "hidden"}
var PrefNameUseGlobalTarget = []string{"useglobaltarget", "useglobaltgt", "use-global", "use-globaltarget", "use_global_target", "use-global-target", "globaltarget"}
var PrefNameGlobalTargetPath = []string{"globaltargetpath", "targetpath", "global_target_path", "global-target-path"}

func (p *prefs) SetByName(name string, val bool) error {
	name = quickclean(name)
	switch {
	case slices.Contains(PrefNameKeepRepo, name):
		tempData.Modify()
		p.KeepRepo = val
	case slices.Contains(PrefNameKeepHidden, name):
		tempData.Modify()
		p.KeepHidden = val
	case slices.Contains(PrefNameUseGlobalTarget, name):
		tempData.Modify()
		p.GlobalTarget = val
	default:
		return ErrBadKey
	}
	return nil
}

func (p *prefs) SetOpt(opt ConfigOption, val bool) {
	switch opt {
	case OptBoolKeepHidden:
		tempData.Modify()
		p.KeepHidden = val
	case OptBoolKeepRepo:
		tempData.Modify()
		p.KeepRepo = val
	case OptBoolUseGlobalTarget:
		tempData.Modify()
		p.GlobalTarget = val
	}
}

func (gd *globalData) SetGlobalTargetPath(path string) { gd.GlobalTargetPath = pops.CleanPath(path) }

func (p *prefs) OverwriteRaw(newp prefs) error {
	if !p.equal(newp) {
		p.KeepHidden = newp.KeepHidden
		p.KeepRepo = newp.KeepRepo
		p.GlobalTarget = newp.GlobalTarget
	}
	return nil
}

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
