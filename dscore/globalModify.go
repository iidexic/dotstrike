package dscore

import (
	"errors"
	"strings"
)

var ErrNotUnique error = errors.New("Attempted to set component alias to a non-unique value")
var ErrParentNotFound error = errors.New("Component.Parent did not match any existing alias.")
var ErrAliasNotFound error = errors.New("Component.Parent did not match any existing alias.")

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

func (gm *globalModify) NewCfg(alias string) *spec {
	s := spec{Ctype: cfgComponent}
	gm.Specs = append(gm.Specs, s)
	sr := &gm.Specs[len(gm.Specs)-1]
	return sr
}
func (G *globals) AddCfg(data spec) {
	G.data.Specs = append(G.data.Specs, data)
}

func (p *prefs) Set(mpref map[string]bool) error {
	for k, b := range mpref {
		switch strings.ToLower(k) {
		case "keeprepo", "keep-repo":
			p.KeepRepo = b
		case "keephidden", "keep-hidden":
			p.KeepHidden = b
		case "globaltarget", "global-target", "global_target":
			p.GlobalTarget = b
		}
	}
	return nil
}

func (pc *pathComponent) SetAlias(alias string) error {
	//TODO: Fix this. should be on cfg level as alias should be unique
	cfptr := gd.data.GetCfg(pc.Parent)
	if cfptr == nil {
		return ErrParentNotFound
	}
	var existingpc *pathComponent
	if pc.Ptype == sourceComponent {
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
