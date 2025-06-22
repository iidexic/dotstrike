package dscore

/*
component: getAlias(), getCtype()
pathComponent:
	Path    string   `toml:"path"`
	Ptype   pathType `toml:"ptype"`
	Ctype   componentType
	Abspath string   `toml:"abspath"`
	Ignores []string `toml:"ignores"`
	Alias   string   `toml:"alias"`
cfg:
	Alias     string          `toml:"alias"`     // name, unique
	Sources   []pathComponent `toml:"sources"`   // files or directories marked as origin points
	Targets   []pathComponent `toml:"targets"`   // files or directories marked as destination points
	Ignorepat []string        `toml:"ignores"`   // ignore patterns that apply to all sources
	Overrides prefs           `toml:"overrides"` //map of settings that will be prioritized over global set
	Ctype     componentType
*/

// ? still uncertain
type BuildComponent struct {
}

// ? Wrap all data requirements into a single function?
// Or take partial cfg data and offload remainder to cfg methods?
func (G *globals) MakeCfg(alias string) (*cfg, error) {

	return &cfg{}, nil //TEMP
}
