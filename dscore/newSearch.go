package dscore

// type pathComponent struct {
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

func FindSpecAlias(aliasP string) (int, string) {
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
		}
		// if aliasP matches 90% of c.alias and length has tolerance of +/- 1 char
		if matchCount >= int(float32(ls)*0.9/float32(ls)) && ls-1 <= len(aliasP) && len(aliasP) <= ls+1 {
			closest = s.Alias
			foundClose = true
		}
	}
	//calculate closest
	if foundClose {
		return -1, closest
	}
	return 0, ""

}
