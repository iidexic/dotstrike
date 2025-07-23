package dscore

import (
	"regexp"
	"strings"
)

type matchMode int

const (
	_ matchMode = iota
	matchExact
	matchExactNCS
	matchPattern
	matchSubstring
	matchFuzzy
	matchCharacters
	matchPercentage
	matchNone
)

//  pathComponent
// Path    string   `toml:"path"`
// Ptype   pathType `toml:"ptype"` //if targetComponent: required to be dirPath
// Ctype   componentType
// Abspath string `toml:"abspath"`
// Ignores []string `toml:"ignores"`
// Alias   string   `toml:"alias"`
// }

// type cfg struct {
// 	Alias     string          `toml:"alias"`     // name, unique
// 	Sources   []pathComponent `toml:"sources"`   // files or directories marked as origin points
// 	Targets   []pathComponent `toml:"targets"`   // files or directories marked as destination points
// 	Ignorepat []string        `toml:"ignores"`   // ignore patterns that apply to all sources
// 	Overrides prefs           `toml:"overrides"` //map of settings that will be prioritized over global set
// 	Ctype     componentType
// }

// type OpMode int
// const (
// 	ModifyComponent OpMode = iota
// 	NewComponent
// 	DeleteComponent
// )
//

// New Approach: maybe just don't put myself in a spot where I have to guess user intent

// SearchConfigs returns an int (tri-state) and string:
//

// FindSpecExact returns a pointer to the cfg object where alias = aliasP
// If aliasP not found, returns nil
func FindSpecExact(aliasP string) *spec {
	for _, s := range gd.data.Specs {
		if aliasP == s.Alias {
			return &s
		}
	}
	return nil

}

// SelectSpec and write selection (index i where gd.Cfgs[i].Alias == alias)
func SelectSpec(alias string) bool {
	for i, s := range gd.data.Specs {
		if alias == s.Alias {
			if !tempData.initialized {
				InitTempData()
			}
			tempData.Selected = i
			e := gd.EncodeIfNeeded(&tempData)
			if e != nil { //TODO: send up for cmd to log?
				panic(e)
			}

			return true
		}
	}
	return false
}

// Not currently using.
func FindSpec(aliasP string) (int, string) {
	var matchCount int
	var foundClose bool
	// maybe just cut the fuzzy match and only run rigid match.
	var close []string
	var closest string
	// I don't know for sure
	_, _ = close, closest
	Global := gd
	for _, s := range Global.data.Specs {
		ls := len(s.Alias)
		if aliasP == s.Alias {
			return 1, aliasP
		}
		// iron out minor spelling mistakes
		for i := range ls - 1 { //TODO: make sure Go is inclusive start-exclusive end. I forget but 90% sure
			if aliasP[i:i+1] == s.Alias[i:i+1] {
				matchCount++

			}
			// if aliasP matches 90% of c.alias and length has tolerance of +/- 1 char
			if matchCount >= int(float32(ls)*0.9/float32(ls)) &&
				ls-1 <= len(aliasP) && len(aliasP) <= ls+1 {
				closest = s.Alias
				foundClose = true
			}
		}
		//calculate closest
		if foundClose {
			return -1, closest
		}
	}
	return 0, ""

}

// hardclean exists. idk
/*
func hardclean(s string) string {
	rmv := []rune(`
!@#$%^&*()_+{}/\|:"'.,;:[]<>?-=~`)
	rmv = append(rmv, '`')
	for _, v := range rmv {
		s = strings.ReplaceAll(s, string(v), "")
	}
	return quickclean(s)
}
*/
func quickclean(s string) string { return strings.TrimSpace(strings.ToLower(s)) }
func checkmatch(lookup string, record string, mode matchMode) int {
	switch mode {
	case matchExact:
		if lookup == record {
			return 1
		}
	case matchExactNCS:
		if quickclean(lookup) == quickclean(record) {
			return 1
		}
	case matchPattern:
		b := make([]byte, len(record))
		copy(b, record)
		m, err := regexp.Match(lookup, b)
		if err != nil {
			return 0
		}
		if m {
			return 1
		}
	case matchSubstring:
		if strings.Contains(quickclean(record), quickclean(lookup)) {
			return 1
		}
	case matchFuzzy:
		lsl := []rune(lookup)
		for i, r := range lsl {
			_, _ = i, r
		}
	case matchPercentage:

	}
	return 0
}

func (g *globalData) FFindSpec(aliasP string) (string, matchMode) {
	for _, s := range g.Specs {
		ls := len(s.Alias)
		if quickclean(aliasP) == quickclean(s.Alias) {
			return aliasP, matchExact
		}
		// iron out minor spelling mistakes
		for i := range ls - 1 { //TODO: make sure Go is inclusive start-exclusive end. I forget but 90% sure
			if aliasP[i:i+1] == s.Alias[i:i+1] {
			}
		}
	}
	return "", 0

}
