package dscore

import (
	"fmt"
	"os"
	"slices"

	"github.com/BurntSushi/toml"
	pops "iidexic.dotstrike/pathops"
)

func ifer(e error) {
	if e != nil {
		panic(e)
	}
}

var ( //TODO: align these errors a bit more with whatever standard is
	ErrNotModified    error = fmt.Errorf("Attempted write of un-modified temp data")
	ErrModifiedNoInit       = fmt.Errorf("Attempted write of modified UN-INITIALIZED temp data")
	ErrNoInit               = fmt.Errorf("Attempted write of un-initialized temp data")
)

//var errNoTemp error = errors.New("TempData is not initialized or does not exist")

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
	// moved TargetPath to Prefs
	//TargetPath string `toml:"storagePath, omitempty"`
	Selected int    `toml:"SelectedSpec"`
	Prefs    prefs  `toml:"prefs"`
	Specs    []spec `toml:"specs, omitempty"`
}

func (g *globalData) equal(g2 *globalData) bool {
	return g.Prefs.equal(g2.Prefs) &&
		g.Selected == g2.Selected &&
		slices.EqualFunc(g.Specs, g2.Specs, specEqual)
}

type globalModify struct {
	*globalData
	initialized, Modified bool
}

/*
	NOTE:FROM DOCS:

|If the "omitempty" option is present the following value will be skipped:
| -> arrays, slices, maps, and string with len of 0, struct with all zero values, bool false
|+ALSO If omitzero is given all int and float types with a value of 0 will be skipped.
| ( start iota at 1 with _ = iota to get past this)
*/

// prefs holds preferences for component-based operations
// used scoped globally or to individual components/parents
type prefs struct {
	KeepRepo         bool   `toml:"keepRepo"`
	KeepHidden       bool   `toml:"keepHidden"`
	GlobalTarget     bool   `toml:"globalTarget.enabled"`
	GlobalTargetPath string `toml:"globalTarget.path"`
	//TODO: symlink handling + symlink preference
}

func (p prefs) equal(p2 prefs) bool {
	return p.KeepHidden == p2.KeepHidden && p.KeepRepo == p2.KeepRepo &&
		p.GlobalTarget == p2.GlobalTarget && p.GlobalTargetPath == p2.GlobalTargetPath
}

// TempGlob exists to store new global data temporarily during runtime
// this will then be checked/merged with GD, and written to globals file
var tempData globalModify

func (G *globals) decodeRawData() {
	md, err := toml.Decode(G.rawContents, &G.data)
	if err != nil {
		panic(fmt.Errorf("Error in dscore DecodeRawData() from data toml\n%e", err))
	}
	G.md = md //? Is this used at all

	//TODO: run CheckDataDecode on debug flag
	//CheckDataDecode(G.data, md)
}

func GetTempData() *globalModify {
	if tempData.initialized {
		return &tempData
	} else {
		return nil
	}
}

func IsDir(ospath string) bool {
	ps, e := os.Stat(ospath)
	if e != nil {
		if os.IsNotExist(e) {
			return false
		}
		return false //TODO: fix function or remove
	}
	if ps.IsDir() {
		return true
	}
	return false
}

// InitTempData populates tempdata from globaldata
// fields populated are required to avoid data loss on toml encode
func InitTempData() {
	if !tempData.initialized {
		tempData = globalModify{
			globalData:  &globalData{},
			initialized: true,
			Modified:    false,
		}
		tempData.Prefs.GlobalTargetPath = gd.data.Prefs.GlobalTargetPath
		tempData.Prefs.GlobalTarget = gd.data.Prefs.GlobalTarget
		tempData.Prefs.KeepHidden = gd.data.Prefs.KeepHidden
		tempData.Prefs.KeepRepo = gd.data.Prefs.KeepRepo
		tempData.Selected = gd.data.Selected
	}
}

// TODO: replace
func (G *globals) EncodeIfNeeded(tg *globalModify) error {
	if tempData.initialized && tempData.Modified {
		return tg.encodeModified()
	} else if tempData.initialized {
		return ErrNotModified
	} else if tempData.Modified {
		return ErrModifiedNoInit

	}
	return ErrNoInit
}

// should only be used when very first writing a non-existent dotstrikeData.toml
func (G *globals) encodeDefaults() error {
	file, e := pops.MakeOpenFileF(globalsFilepath())
	if e != nil {
		return e
	}
	defer file.Close()
	encode := toml.NewEncoder(file)
	e = encode.Encode(G.data)
	if e != nil {
		return e
	} else {
		return nil
	}
}

// encodeModified gm data exclusively to main toml
func (gm *globalModify) encodeModified() error {
	file, e := pops.OpenFileRW(globalsFilepath())
	if e != nil || file == nil {
		return e
	}
	defer file.Close()
	encode := toml.NewEncoder(file)
	e = encode.Encode(gm.globalData)
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
