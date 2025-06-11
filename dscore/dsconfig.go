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

// globals holds configuration status and data
// globals must be read from file in config step every time ds is run
type globals struct {
	data          globalData
	status        globalsReadResult
	loaded        bool
	checkedpaths  []string
	rawContents   string
	GlobalMessage []string
}
type globalData struct {
	//check later if omitempty is needed
	cfgs         []cfg `toml:"cfgs, omitempty"`
	prefs        `toml:"prefs, omitempty"`
	dsconfigPath string `toml:"targetPath,omitempty"`
}

type prefs struct {
	keepRepo     bool
	keepHidden   bool
	globalTarget bool
}

type allTemp struct {
	components []component
	aliases    []string
}

// tempGlob exists to store new global data temporarily during runtime
// tempGlob will then be checked + merged with main globals, and written to globals file
var tempGlob globals

// tempComponent stores changes to a component before being merged with that component and written to globals file
// may be superfluous; can use tempGlob.cfgs for this
var tempComponent allTemp

// func (G *globals) exists() { }// Uncertain of original intent. Most likely covered by G.Getcomponent

func (G *globals) loadFromRaw() {
	fromTomlString(G.rawContents)
}

func (G *globals) output(outStr string) {
	G.GlobalMessage = append(G.GlobalMessage, outStr)
}
func (G *globals) outputf(outStr string, anyfmt ...any) {
	G.GlobalMessage = append(G.GlobalMessage, fmt.Sprintf(outStr, anyfmt...))
}

func (G *globals) Dump() []string {
	dump := []string{
		"__GLOBALS__",
		G.status.string(),
		fmt.Sprintf("globals loaded: %t", G.loaded),
		fmt.Sprintf("user cfgs: %v", G.data.cfgs),
		fmt.Sprintf("preferences: %+v", G.data.prefs),
		fmt.Sprintf("globals file path: %s", G.data.dsconfigPath),
		fmt.Sprintf("checked paths: %v", G.checkedpaths),
		"__MESSAGES__",
	}
	dump = append(dump, G.GlobalMessage...)
	return dump
}

func (G globals) DumpRaw() string {
	return fmt.Sprintf("%+v", G)
}

func (G globals) DecodeRawData() {
	toml.Decode(G.rawContents, G.data)

}
