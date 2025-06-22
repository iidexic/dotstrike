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

/*
	NOTE:FROM DOC:

|If the "omitempty" option is present the following value will be skipped:
|-> arrays, slices, maps, and string with len of 0
|-> struct with all zero values
|-> bool false
|+ALSO If omitzero is given all int and float types with a value of 0 will be skipped.
*/
type prefs struct {
	KeepRepo     bool `toml:"keepRepo"`
	KeepHidden   bool `toml:"keepHidden"`
	GlobalTarget bool `toml:"globalTarget"`
}
type modifyData struct {
	*globalData
	modified bool
}

// TempGlob exists to store new global data temporarily during runtime
// this will then be checked/merged with GD, and written to globals file
var TempGlob *globals = nil
var TempData *globalData = nil
var DataModified bool = false

func IsTempData() bool { return TempGlob != nil }

func GetTempGlobals() *globals {
	// dont leave this like this pleaze
	if TempGlob == nil {
		TempGlob = &globals{
			data: globalData{},
		}
	}

	return TempGlob
}

// NOTE: globalData is the only piece that requires modification when making changes
func getTempData() *globalData {
	// dont leave this like this pleaze
	if TempData == nil {
		TempData = &globalData{
			Prefs: GD.data.Prefs,
		}
	}

	return TempData
}

func (G *globals) decodeRawData() {
	md, err := toml.Decode(G.rawContents, &G.data)
	if err != nil {
		panic(fmt.Errorf("Error in dscore DecodeRawData() from data toml\n%e", err))
	}
	G.md = md //? Is this used at all

	//TODO: run CheckDataDecode on debug flag
	//CheckDataDecode(G.data, md)
}

func (G *globals) EncodeIfNeeded(tempg globals) {
	if DataModified {

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
