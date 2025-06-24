package dscore

import (
	"fmt"
	"slices"

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
	Cfgs       []cfg  `toml:"cfgs, omitempty"`
	Prefs      prefs  `toml:"prefs"`
	TargetPath string `toml:"storagePath, omitempty"`
	Selected   int    `toml:"SelectedCFG"`
}

func (g *globalData) equal(g2 *globalData) bool {
	return g.Prefs.Equal(g2.Prefs) &&
		g.TargetPath == g2.TargetPath &&
		g.Selected == g2.Selected &&
		slices.EqualFunc(g.Cfgs, g2.Cfgs, cfgEqual)

}
func (p prefs) Equal(p2 prefs) bool {
	return p.KeepHidden == p2.KeepHidden && p.KeepRepo == p2.KeepRepo &&
		p.GlobalTarget == p2.GlobalTarget
}

/*
	NOTE:FROM DOCS:

|If the "omitempty" option is present the following value will be skipped:
| -> arrays, slices, maps, and string with len of 0, struct with all zero values, bool false
|+ALSO If omitzero is given all int and float types with a value of 0 will be skipped.
| ( start iota at 1 with _ = iota to get past this)
*/
type prefs struct {
	KeepRepo     bool `toml:"keepRepo"`
	KeepHidden   bool `toml:"keepHidden"`
	GlobalTarget bool `toml:"globalTarget"`
}

type globalModify struct {
	*globalData
	initialized, Modified bool
}

// TempGlob exists to store new global data temporarily during runtime
// this will then be checked/merged with GD, and written to globals file
var TempData globalModify

func (G *globals) decodeRawData() {
	md, err := toml.Decode(G.rawContents, &G.data)
	if err != nil {
		panic(fmt.Errorf("Error in dscore DecodeRawData() from data toml\n%e", err))
	}
	G.md = md //? Is this used at all

	//TODO: run CheckDataDecode on debug flag
	//CheckDataDecode(G.data, md)
}

func InitTempData() {
	if !TempData.initialized {
		TempData = globalModify{
			globalData:  &globalData{},
			initialized: true,
			Modified:    false,
		}
		TempData.Prefs.GlobalTarget = GD.data.Prefs.GlobalTarget
		TempData.Prefs.KeepHidden = GD.data.Prefs.KeepHidden
		TempData.Prefs.KeepRepo = GD.data.Prefs.KeepRepo
		TempData.TargetPath = GD.data.TargetPath
		TempData.Selected = GD.data.Selected
	}
}

// TODO: replace
func (G *globals) EncodeIfNeeded(tg globalModify) {
	if TempData.Modified {

	}
}

// encodeG is functional encode
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

func (G *globals) logG(outStr string) {
	G.GlobalMessage = append(G.GlobalMessage, outStr)
}
func (G *globals) logfG(outStr string, anyfmt ...any) {
	G.GlobalMessage = append(G.GlobalMessage, fmt.Sprintf(outStr, anyfmt...))
}
