package dscore

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

//var symbols = []rune(` !@#$%^&*()_+{}/\|:"'.,;:[]<>?-=~`)

//TODO:(low) Identify what is being used instead. Then Delete or Develop FindSpec

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
