package dscore

import (
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

var symbols = []rune(` !@#$%^&*()_+{}/\|:"'.,;:[]<>?-=~`)

func FindSpecExact(aliasP string) *Spec {
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

// performs an exact search and a likeness search.
// int returned is 1 if exact match found, -1 if there is a 90% likeness, and 0 otherwise
// string returned is matching spec's alias; otherwise returns 0, ""
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
		for i := range ls - 1 {
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

// fuzzysearch checks whether text contains each character in substring, in the same relative character position.
//
// returns true if:
//  1. all runes in sub exist within text
//  2. each rune @ position n in sub exists at some position i in text such that i(rune[n]) > i(rune[n-1])
//
// ex: fuzzySearch("afbcdef", "fdf",false) returns `true`
func fuzzySearch(text string, sub string, caseSens bool) bool {
	if len(text) < len(sub) {
		return false
	}
	if !caseSens {
		text = strings.ToLower(text)
		sub = strings.ToLower(sub)
	}
	subslice := []rune(sub)

	// Pure Text Method
	for _, l := range subslice {
		runepos := strings.IndexRune(text, l)
		if runepos < 0 {
			return false
		}
		lastIndex := len(text) - 1
		if lastIndex > runepos {
			text = text[runepos+1:]
		} else {
			text = ""
		}
	}
	return true
}

func stripSymbols(s string) string {
	rmv := []rune(`
!@#$%^&*()_+{}/\|:"'.,;:[]<>?-=~`)
	rmv = append(rmv, '`')
	for _, v := range rmv {
		s = strings.ReplaceAll(s, string(v), "")
	}
	return QuickClean(s)
}
func hasSymbols(s string) bool { return strings.ContainsAny(s, string(symbols)) }

// QuickClean performs some string standardization to improve matching and lookups for dscore functions
func QuickClean(s string) string { return strings.TrimSpace(strings.ToLower(s)) }

/* func checkmatch(lookup string, record string, mode matchMode) int {
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
} */

func (g *globalData) FFindSpec(aliasP string) (string, matchMode) {
	for _, s := range g.Specs {
		ls := len(s.Alias)
		if QuickClean(aliasP) == QuickClean(s.Alias) {
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
