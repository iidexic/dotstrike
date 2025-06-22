package dscore

import (
	"slices"
	"strings"
)

type FindType int

const (
	FindNil FindType = iota
	FindAny
	FindCfg
	FindPathComp
	FindSource
	FindTarget
	BoundPathComp
	BoundSource
	BoundTarget
)

// Lookup contains a generated or manually passed list of search options
// as well as containing the results of a performed search
// When a search is run, FindType is transformed into a Lookup var
// used to determine what to search and return
type search struct {
	GetType componentType
	aliases []string
	bind    searchResult
	output  searchResult
}
type searchResult interface {
}

type Lookup struct {
	//configs, sources, targets, search bound to (selected config object??)
	GetCfg, GetSrc,
	GetTgt, BoundOnly bool
	findCfgs, findSources, findTargets []string
	foundCfgs                          []cfg
	foundSources, foundTargets         []pathComponent
	// MatchFound returns true if rigid search returned result, and false otherwise
	cfgMatchFound, srcMatchFound, tgtMatchFound bool
}

func (L Lookup) componentTypes() []componentType {
	c := []componentType{}
	if L.GetCfg {
		c = append(c, cfgComponent)
	}
	if L.GetSrc {
		c = append(c, sourceComponent)
	}
	if L.GetTgt {
		c = append(c, targetComponent)
	}
	return c

}

// GetBoundComponents finds pathComponents matching/containing aliasPattern within parent cfg
func GetBoundComponents(parent *cfg, aliasPatterns []string, request FindType) {
	// is anything being modified here?
}

// CfgData returns cfg with exact match alias
func (G globals) CfgData(alias string) Lookup {
	// NOTE: started as a globalData method. keep most/all
	// methods under globals until a reason to separate comes up
	l := Lookup{GetCfg: true, foundCfgs: make([]cfg, 1), cfgMatchFound: false}
	for _, c := range G.data.Cfgs {
		if c.Alias == alias {
			l.foundCfgs[0] = c
			l.cfgMatchFound = true
			break
		}
	}
	return l
	/*TODO: Implement usable return: Using Lookup
	- return empty cfg | error return */
	//TODO: determine which globals methods should be pointer methods
}

// NOTE: started as a globalData method. keep most/all
// methods under globals until a reason to separate comes up
func (G globals) SourceData(parent cfg, alias string) Lookup {
	l := Lookup{GetSrc: true, foundSources: make([]pathComponent, 1), srcMatchFound: false}
	for _, c := range G.data.Cfgs {
		if c.Alias == alias {
			l.foundCfgs[0] = c
			l.cfgMatchFound = true
			break
		}
	}
	return l
	/*TODO: Implement usable return: Using Lookup
	- return empty cfg | error return */
	//TODO: determine which globals methods should be pointer methods
}

// GetComponents performs search for user data based on full/partial alias and FindType provided
func (G globals) GetComponents(aliasPattern []string, request FindType) []component {
	look := Lookup{
		GetCfg:    slices.Contains([]FindType{FindAny, FindCfg}, request),
		GetSrc:    slices.Contains([]FindType{FindAny, FindPathComp, FindSource, BoundPathComp, BoundSource}, request),
		GetTgt:    slices.Contains([]FindType{FindAny, FindPathComp, FindTarget, BoundPathComp, BoundTarget}, request),
		BoundOnly: slices.Contains([]FindType{BoundPathComp, BoundSource, BoundTarget}, request), //TODO: Pull into GetBoundComponents?
	}
	_ = look
	searchResult := []component{}
	/* Return Data:
	this one returns components. Build out Interface a bit? Or another struct with just a bunch of slices to slap shit in
	Or could do []ints that store index within Global */
	return searchResult
}

// TODO: Finish Building Search!
// find_cfg searches all existing cfg aliases for pattern:
func (gd globalData) find_cfg(pattern string) (bool, []*cfg) {

	clist := make([]*cfg, 0, len(gd.Cfgs))
	gotMatch := true

	casesens := func(c cfg) string { return c.Alias }
	if strings.ToLower(pattern) == pattern {
		casesens = func(c cfg) string { return strings.ToLower(c.Alias) }
	}
	for _, cf := range gd.Cfgs {
		dcomp := strings.Contains(casesens(cf), pattern)
		if dcomp {
			clist = append(clist, &cf)
		}
	}

	// if nothing was found, find if has any same chars
	if len(clist) == 0 {
		gotMatch = false
		for _, cf := range gd.Cfgs {
			rc := 0.0
			findIn := casesens(cf)
			for _, r := range pattern {
				if strings.ContainsRune(findIn, r) {
					rc += 1.0
				}
			}
			// completely arbitrary condition
			if rc >= 0.4*float64(len(findIn)) {
				clist = append(clist, &cf)
			}
		}
	}
	return gotMatch, clist
}
