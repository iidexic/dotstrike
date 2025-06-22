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

// FindCfgExact returns a pointer to the cfg object where alias = aliasP
// If aliasP not found, returns nil
func FindCfgExact(aliasP string) *cfg {
	for _, c := range GD.data.Cfgs {
		if aliasP == c.Alias {
			return &c
		}
	}
	return nil

}

// SelectCfg and write to temp so it can be encoded
func SelectCfg(aliasP string) bool {
	for i, c := range GD.data.Cfgs {
		if aliasP == c.Alias {

			TempGlob.data.Selected = i
			GD.EncodeIfNeeded(*TempGlob)
			return true
		}
	}
	return false
}

func FindCfgAlias(aliasP string) (int, string) {
	var matchCount int
	var foundClose bool
	// maybe just cut the fuzzy match and only run rigid match.
	var close []string
	var closest string
	// I don't know for sure
	_, _ = close, closest
	Global := GD
	for _, c := range Global.data.Cfgs {
		lc := len(c.Alias)
		if aliasP == c.Alias {
			return 1, aliasP
		}
		// iron out minor spelling mistakes
		for i := range lc - 1 { //TODO: make sure Go is inclusive start-exclusive end. I forget but 90% sure
			if aliasP[i:i+1] == c.Alias[i:i+1] {
				matchCount++
			}
		}
		// if aliasP matches 90% of c.alias and length has tolerance of +/- 1 char
		if matchCount >= int(float32(lc)*0.9/float32(lc)) && lc-1 <= len(aliasP) && len(aliasP) <= lc+1 {
			closest = c.Alias
			foundClose = true
		}
	}
	//calculate closest
	if foundClose {
		return -1, closest
	}
	return 0, ""

}
