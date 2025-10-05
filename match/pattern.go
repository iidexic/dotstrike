package match

import (
	"fmt"
	"strings"
)

type PreIgnoreList []string

type PreIgnoreData struct {
	*PreIgnoreList
}

// Textpattern takes comparison strings as input and output the result of the match (boolean Matches/Doesn't)
// They all have a function to Set the pattern that will be matched against.
type TextPattern interface {
	Matches(string) bool
	Set(string, bool) bool
	IsSet() bool
}

var (
	ErrEmptyPattern = fmt.Errorf("No Pattern provided when making TextPattern (Nothing to Match with)")
)

// utility only
type truePtn struct{}

func (truePtn) Matches(x string) bool         { return true }
func (truePtn) Set(na string, none bool) bool { return false }
func (truePtn) IsSet() bool                   { return true }

// PathPattern is a quirky little guy :)
type PathPattern struct {
	baseptn   TextPattern
	parentptn TextPattern
	rootptn   TextPattern
	local     bool
}

// TODO:(hi) Uh test/document all this
func (pp *PathPattern) Set(ptn string, overwrite bool) bool {
	switch {
	case pp.IsSet() && !overwrite:
		return false
	case strings.HasPrefix(ptn, "./") || strings.HasPrefix(ptn, ".\\"):
		if strings.Count(ptn, "/") > 1 || strings.Count(ptn, "\\") > 1 {
			pp.local = true
			ptn = ptn[2:]
			i := max(strings.LastIndex(ptn, "/"), strings.LastIndex(ptn, "\\"))
			if i < 0 { //technically impossible
				return false
			}
			base := ptn[i+1:]
			parent := ptn[:i]
			pp.baseptn, pp.parentptn = &SubPattern{}, &SubPattern{}
			pp.baseptn.Set(base, true)
			pp.parentptn.Set(parent, true)
			break
		}
		fallthrough
	case strings.HasPrefix(ptn, "*/") || strings.HasPrefix(ptn, "*\\"):
		pp.baseptn = &SubPattern{}
		b := pp.baseptn.Set(ptn[2:], true)
		if b {
			pp.rootptn, pp.parentptn = truePtn{}, truePtn{}
			return true
		} else {
			return false
		}
	}

	return false
}

func (pp PathPattern) IsSet() bool {
	return pp.baseptn.IsSet() && pp.parentptn.IsSet() && pp.rootptn.IsSet()
}
func (pp PathPattern) Matches(input string) bool {
	return pp.baseptn.Matches(input) && pp.parentptn.Matches(input) && pp.rootptn.Matches(input)
}

// NewSubptn is a helper function to generate a preset config subpattern.
//   - `matchAnywhere = true` makes a pure subpattern search
func NewSubptn(ptn string, matchAnywhere bool) TextPattern {
	var pre, suf bool
	// kinda half-assed, whatever
	if ptn == "*" || ptn == "**" {
		return &truePtn{}
	}
	if lp := len(ptn); lp > 0 {
		if lp > 1 {
			ptn, pre, suf = trimwild(ptn)
		} else {
			pre, suf = true, true
		}
		if matchAnywhere {
			return &SubPattern{string: ptn}
		}
		return &SubPattern{string: ptn, Prefix: pre, Suffix: suf}
	}
	return nil
}

// trimwild expects a NON-EMPTY STRING
// Also it is way too late I am just going to assume this is fine
// WARNING: THIS IS PROBABLY NOT FINE
func trimwild(ptn string) (string, bool, bool) {
	var mpre, msuf bool = true, true
	if ptn[0] == '*' && len(ptn) > 1 {
		mpre = false
		ptn = ptn[1:]
	}
	if lp := len(ptn); ptn[lp-1] == '*' && lp > 1 {
		msuf = false
		ptn = ptn[:len(ptn)-1]
	}
	return ptn, mpre, msuf
}

// SubPattern matches if the input string contains the set pattern in the configuration determined by SubPattern's booleans.
//   - If !Prefix && !Suffix (default), the pattern can appear anywhere in the input
//   - With Prefix or Suffix set, matches if pattern appears at the very beginning or very end of the input string, respectively.
//   - With Prefix && Suffix, matches if pattern appears at both beginning and end. These do not have to be the same occurrence.
type SubPattern struct {
	string
	Prefix, Suffix bool
	set            bool
}

func (sp *SubPattern) Set(ptn string, overwrite bool) bool {
	if sp.set && !overwrite {
		return false
	}
	if !strings.HasPrefix(ptn, "*") {
		sp.Prefix = true
		ptn = ptn[1:]
	}
	if !strings.HasSuffix(ptn, "*") {
		sp.Suffix = true
		ptn = ptn[:len(ptn)-1]
	}
	sp.string = ptn
	sp.set = true
	return true
}
func (sp SubPattern) IsSet() bool { return sp.set }

func (sp SubPattern) Matches(input string) bool {
	switch {
	case sp.Prefix && sp.Suffix:
		return strings.HasPrefix(input, sp.string) && strings.HasSuffix(input, sp.string)
	case sp.Prefix:
		return strings.HasPrefix(input, sp.string)
	case sp.Suffix:
		return strings.HasSuffix(input, sp.string)
	}
	return strings.Contains(input, sp.string)
}
