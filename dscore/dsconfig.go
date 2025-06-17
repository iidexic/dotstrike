package dscore

import (
	"fmt"

	"github.com/BurntSushi/toml"
	pops "iidexic.dotstrike/pathops"
)

func ifer(e error) {
	if e != nil {
		panic(e)
	}
}

// intended create writer; just using opened file
// type stwriter struct { stringout string }
// func (s stwriter) Write(b []byte) (n int, e error) { }

// ╭─────────────────────────────────────────────────────────╮
// │              Dotstrike Config+Data Structs              │
// ╰─────────────────────────────────────────────────────────╯

// globals holds configuration status and data
// globals must be read from file in config step every time ds is run

type globals struct {
	data          globalData //May need pointer to make writeable (if needed). Pointer may cause toml decode issue
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
	Cfgs       []cfg  `toml:"cfgs, omitempty"`
	Prefs      prefs  `toml:"prefs, omitempty"`
	TargetPath string `toml:"storagePath"`
	Selected   int    `toml:"SelectedCFG"`
}

type prefs struct {
	KeepRepo     bool `toml:"keepRepo, omitempty"`
	KeepHidden   bool `toml:"keepHidden, omitempty"`
	GlobalTarget bool `toml:"globalTarget, omitempty"`
}

// TempGlob exists to store new global data temporarily during runtime
// this will then be checked/merged with GD, and written to globals file
var TempGlob globals

func (G *globals) decodeRawData() {
	md, err := toml.Decode(G.rawContents, &G.data)
	if err != nil {
		panic(fmt.Errorf("Error in dscore DecodeRawData() from data toml\n%e", err))
	}
	G.md = md //? Is this used at all

	//TODO: run CheckDataDecode on debug flag
	//CheckDataDecode(G.data, md)
}

func (G *globals) encodeG() error {
	testpath := ""
	file := pops.MakeOpenFileF(testpath)
	defer file.Close()
	encode := toml.NewEncoder(file)
	e := encode.Encode(G.data)
	if e != nil {
		return e
	} else {
		return nil
	}
}

func (G *globals) loadFromRaw() {
	fromTomlString(G.rawContents)
}

func (G *globals) logG(outStr string) {
	G.GlobalMessage = append(G.GlobalMessage, outStr)
}
func (G *globals) logfG(outStr string, anyfmt ...any) {
	G.GlobalMessage = append(G.GlobalMessage, fmt.Sprintf(outStr, anyfmt...))
}
