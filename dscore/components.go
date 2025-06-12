package dscore

import "slices"

// Denote whether paths in pathObjects are path or dir
type pathType int

// Denote whether component is source or target. Uncertain of implementation
type componentType int

const (
	filePath pathType = iota
	dirPath
	sourceComponent componentType = iota
	targetComponent
)

// component interface for interop search
type component interface {
	name() string
}

// pathComponent is the core of a source or target;
// contains path info
type pathComponent struct {
	path    string
	ptype   pathType
	ctype   componentType
	abspath string
	ignores []string // local ignores; specific to this dir
	alias   string
}

func (pc pathComponent) name() string { return pc.alias }

// cfg is the primary structure used to define a move/strike
type cfg struct {
	Alias     string            // name, unique
	Sources   []pathComponent   // files or directories marked as origin points
	Targets   []pathComponent   // files or directories marked as destination points
	Ignorepat []string          // ignore patterns that apply to all sources
	Overrides map[string]string //map of settings that will be prioritized over global set
}

func (cc cfg) name() string { return cc.Alias }

type FindType int

const (
	FindAny FindType = iota
	FindCfg
	FindPathComp
	FindSource
	FindTarget
	BoundPathComp
	BoundSource
	BoundTarget
)

// lookup contains a list of options for a search.
// When a search is run, FindType is transformed into a lookup var
// used to determine what to search and return
type lookup struct {
	//configs, sources, targets, search bound to (selected config object??)
	getCfg, getSrc, getTgt, boundOnly bool
}

// Getcomponent performs search for user data based on full/partial alias and FindType provided
func (G *globals) Getcomponent(aliases []string, request FindType) {
	look := lookup{
		getCfg:    slices.Contains([]FindType{FindAny, FindCfg}, request),
		getSrc:    slices.Contains([]FindType{FindAny, FindPathComp, FindSource, BoundPathComp, BoundSource}, request),
		getTgt:    slices.Contains([]FindType{FindAny, FindPathComp, FindTarget, BoundPathComp, BoundTarget}, request),
		boundOnly: slices.Contains([]FindType{BoundPathComp, BoundSource, BoundTarget}, request),
	}
	/*
		Return Data:
		this one returns components. Build out Interface a bit? Or another struct with just a bunch of slices to slap shit in
		Or could do []ints that store index within Global

	*/
	if look.getCfg {

	}
	if look.getSrc {
		if look.boundOnly {

		} else {

		}
	}
	if look.getTgt {
		if look.boundOnly {

		} else {

		}

	}

}
