package dscore

import (
	"errors"
	"strings"
)

var errorNotUnique error = errors.New("Attempted to set component alias to a non-unique value")
var errorParentNotFound error = errors.New("Component.Parent did not match any existing alias.")

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

func (gm *globalModify) NewCfg(alias string) *cfg { return &cfg{Ctype: cfgComponent} }

func (G *globals) AddCfg(data cfg) {
	G.data.Cfgs = append(G.data.Cfgs, data)
}

// SetAlias is not following methodology of other functions, remove
/* func (cc *cfg) SetAlias(alias string) error { if gd.data.GetCfg(alias) != nil { return errorNotUnique } cc.Alias = alias return nil } */
func (cc *cfg) GetIgnores() *[]string { return &cc.Ignorepat }
func (cc *cfg) GetLocalPrefs() *prefs { return &cc.Overrides }

// func (cc *cfg) getChildByPath(path string) *pathComponent { }
func (cc cfg) IsPathChild(path string) bool {
	for _, src := range cc.Sources {
		_ = src
	}
	for _, tgt := range cc.Targets {
		_ = tgt
	}
	return true //NOTE:PLACEHOLDER
}
func (cc *cfg) AddIgnores(ignores []string) {
	cc.Ignorepat = append(cc.Ignorepat, ignores...)
}
func (cc *cfg) CheckAddPath(path string, isSource bool) bool {
	var ctyp componentType
	if isSource {
		ctyp = sourceComponent
	} else {
		ctyp = targetComponent
	}
	if !cc.IsPathChild(path) {
		cc.Sources = append(cc.Sources, *newPathComponent(path, ctyp))
		return true
	}
	return false
}

// CheckAddMultiplePaths adds paths to cfg.Sources if isSource, cfg.Targets if !isSource
func (cc *cfg) CheckAddMultiplePaths(paths []string, isSource bool) {
	var ctyp componentType
	if isSource {
		ctyp = sourceComponent
	} else {
		ctyp = targetComponent
	}
	_ = ctyp
	for _, p := range paths {
		_ = p //if pc := cc.getChildByPath(p); pc == nil { } else { }
	}
}
func (p *prefs) Set(mpref map[string]bool) {
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
}

func (pc *pathComponent) SetAlias(alias string) error {
	//TODO: Fix this
	cfptr := gd.data.GetCfg(pc.Parent)
	if cfptr == nil {
		return errorParentNotFound
	}
	var existingpc *pathComponent
	if pc.Ptype == sourceComponent {
		existingpc = cfptr.getSource(alias)
	} else {
		existingpc = cfptr.getTarget(alias)
	}
	if existingpc != nil {
		return errorNotUnique
	}
	pc.Alias = alias
	return nil
}
