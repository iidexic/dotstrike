package dscore

import (
	"errors"
	"fmt"
	"strings"

	pops "iidexic.dotstrike/pathops"
)

var ErrNotUnique error = errors.New("Attempted to set component alias to a non-unique value")
var ErrParentNotFound error = errors.New("Component.Parent did not match any existing alias.")
var ErrAliasNotFound error = errors.New("Component.Parent did not match any existing alias.")
var ErrBadKey error = errors.New("Provided map key does not exist")

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
// WARN: sets gm.Modified true
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

// Set is option 1 for modifying preferences
func (p *prefs) Set(mpref map[string]bool) error {
	for k, b := range mpref {
		ck := quickclean(k)
		switch ck {
		case "keeprepo", "keep-repo", "keep_repo", "repo":
			p.KeepRepo = b
		case "keephidden", "keep-hidden", "keep_hidden", "hidden":
			p.KeepHidden = b
		case "globaltarget", "globaltgt", "global_target", "global-target":
			p.GlobalTarget = b
		default:
			return ErrBadKey
		}
	}
	return nil
}
func (p *prefs) SetGlobalTargetPath(path string) { p.GlobalTargetPath = pops.CleanPath(path) }
func (p *prefs) OverwriteRaw(newp prefs) error {
	if !p.equal(newp) {
		p.KeepHidden = newp.KeepHidden
		p.KeepRepo = newp.KeepRepo
		p.GlobalTarget = newp.GlobalTarget
		if newp.GlobalTargetPath != "" {
			p.GlobalTargetPath = newp.GlobalTargetPath
		}

	}
	return nil
}

// SetAlias sets the PathComponent alias.
// If PathComponent is not unique, alias is not set, and ErrNotUnique is returned
func (pc *pathComponent) SetAlias(alias string) error {
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
