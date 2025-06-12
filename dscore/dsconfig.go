package dscore

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func ifer(e error) {
	if e != nil {
		panic(e)
	}
}

// ╭─────────────────────────────────────────────────────────╮
// │              Dotstrike Config+Data Structs              │
// ╰─────────────────────────────────────────────────────────╯

// globals holds configuration status and data
// globals must be read from file in config step every time ds is run
type globals struct {
	data          globalData
	status        globalsReadResult
	loaded        bool
	checkedpaths  []string
	rawContents   string
	dsconfigPath  string
	GlobalMessage []string
	md            toml.MetaData
}
type globalData struct {
	//check later if omitempty is needed
	cfgs       []cfg
	Prefs      prefs     `toml:"prefs, omitempty"`
	TargetPath string    `toml:"storagePath"`
	CfgToml    []cfgMake `toml:"cfgs, omitempty"`
}

type prefs struct {
	keepRepo     bool
	keepHidden   bool
	globalTarget bool
}
type cfgMake struct {
	Alias     string            `toml:"alias"`   // name, unique
	Sources   []string          `toml:"sources"` // files or directories marked as origin points
	Targets   []string          `toml:"targets"` // files or directories marked as destination points
	Ignorepat []string          // ignore patterns that apply to all sources
	Overrides map[string]string //map of settings that will be prioritized over global set
}

// TempGlob exists to store new global data temporarily during runtime
// this will then be checked/merged with GD, and written to globals file
var TempGlob globals

func (G globals) DecodeRawData() {
	md, err := toml.Decode(G.rawContents, &G.data)
	if err != nil {
		panic(fmt.Errorf("Error in dscore DecodeRawData() from data toml\n%e", err))
	}
	CheckDataDecode(G.data, md)
}

func (G *globals) loadFromRaw() {
	fromTomlString(G.rawContents)
}

func (G *globals) output(outStr string) {
	G.GlobalMessage = append(G.GlobalMessage, outStr)
}
func (G *globals) outputf(outStr string, anyfmt ...any) {
	G.GlobalMessage = append(G.GlobalMessage, fmt.Sprintf(outStr, anyfmt...))
}
