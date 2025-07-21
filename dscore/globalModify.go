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

func (gm *globalModify) NewSpec(alias string) *spec {
	s := spec{Alias: alias, Ctype: cfgComponent}
	gm.Specs = append(gm.Specs, s)
	sr := &gm.Specs[len(gm.Specs)-1] //works
	gm.Modified = true
	return sr
}

// GetModifiableSpec returns a pointer to a spec that will be encoded when the program exits
// If spec exists but has not yet been modified, this adds that spec to gm.Specs
// WARN: sets gm.Modified true
func (gm *globalModify) GetModifiableSpec(alias string) (*spec, error) {
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
		case "globaltarget", "global target", "global_target", "global-target":
			p.GlobalTarget = b
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
	//TODO: Fix this. should be on cfg level as alias should be unique
	cfptr := gd.data.GetSpec(pc.Parent)
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
