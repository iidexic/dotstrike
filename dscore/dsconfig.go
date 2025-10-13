package dscore

import (
	"fmt"
	"maps"
	"strings"

	"github.com/BurntSushi/toml"
	"iidexic.dotstrike/config"
	"iidexic.dotstrike/uout"
)

// ── Dotstrike Config+Data Structs ───────────────────────────────────

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
	Selected         int    `toml:"SelectedSpec"`
	GlobalTargetPath string `toml:"targetpath"`
	Prefs            prefs  `toml:"prefs"`
	Specs            []Spec `toml:"specs"`
	//TODO: add []rawComponent/implement rawComponent
}

//TODO:(low - general feature) symlink handling + symlink preference

type globalModify struct {
	*globalData
	initialized, Modified bool
}

type prefs struct {
	Bools map[ConfigOption]bool
}

// ── Config Options ──────────────────────────────────────────────────

// TODO:(mid) Replace entirely with config.OptionKey
//
//	Use config.OptionKey directly where needed and then delete this
type ConfigOption = config.OptionKey

const (
	BoolIgnoreRepo   = config.BoolIgnoreRepo
	BoolIgnoreHidden = config.BoolIgnoreHidden
	BoolRootSubdir   = config.BoolRootSubdir
	BoolNoFiles      = config.BoolNoFiles
	//BoolSeparateSources    = config.BoolSourceSubdirs
	BoolCopyAllDirs        = config.BoolCopyAllDirs
	BoolUseGlobalTarget    = config.BoolUseGlobalTarget
	BoolKillGlobalTarget   = config.BoolKillGlobalTarget
	BoolOverrideOn         = config.BoolOverrideOn
	StringGlobalTargetPath = config.StringGlobalTargetPath
	NotAnOption            = config.NotAnOption
)

// Delete if not in use (or probably even if it is)
var OptionID = config.LookupOption

func OptionIsBool(opt ConfigOption) bool   { return config.AllOptions[opt].Type == config.Tbool }
func OptionIsString(opt ConfigOption) bool { return config.AllOptions[opt].Type == config.Tstring }

var GetOption = config.OptFrom

// ──────────────────────────────────────────────────────────────────────

func (G *globals) Detail() string {
	lines := make([]string, 1, 32) //arbitrary
	lines[0] = fmt.Sprintf(`GLOBAL USER DATA:
==================
Config Path: '%s'
Selected Spec(index): %d

Globals Log (instance):
`, G.dsconfigPath, G.data.Selected)
	lines = append(lines, G.GlobalMessage...)

	lines = append(lines, G.data.Prefs.Detail())
	for i := range G.data.Specs {
		lines = append(lines, G.data.Specs[i].Detail())
	}

	return strings.Join(lines, "\n")
}

func (gd *globalData) Detail(verbose bool) string {
	if !verbose {
		return gd.DetailSimple()
	}
	out := uout.NewOut("[ User Data ]")
	out.V("Config:")
	out.IndR()
	out.NV("Global Target Path", gd.GlobalTargetPath)
	out.LV(gd.Prefs.Bools)
	out.IndL()
	out.V("Specs:")
	out.IndR()
	out.LV(gd.Specs)
	return out.String()
}

func (gd *globalData) DetailSimple() string {

	det := make([]string, len(gd.Specs)+len(gd.Prefs.Bools))
	det[0] = "[ User Data ]"
	i := 1
	for n, s := range gd.Specs {
		det[i] = fmt.Sprintf("[%d] Spec %s: %d sources, %d targets, %d overrides", n, s.Alias, len(s.Sources), len(s.Targets), len(s.Overrides.Bools))
		i++
	}
	det[i] = fmt.Sprintf("Prefs: %d boolean prefs", len(gd.Prefs.Bools))
	return strings.Join(det, "\n")
}

func (p prefs) Detail() string {
	if len(p.Bools) == 0 {
		return "0 options set"
	}
	out := fmt.Sprintf("%d options set", len(p.Bools))
	for k, v := range p.Bools {
		out = fmt.Sprintf("%s\n%s:%t", out, k.String(), v)
	}
	return out
}

func (p prefs) equal(p2 prefs) bool {
	return maps.Equal(p.Bools, p2.Bools)
}

// tempData is the location where ALL changes to user data are written to before encode
var tempData globalModify

func (G *globals) decodeRawData() {
	md, err := toml.Decode(G.rawContents, &G.data)
	if err != nil {
		panic(fmt.Errorf("Error in dscore DecodeRawData() from data toml\n%w", err))
	}
	//TODO:(VO.1) REMOVE UNUSED
	G.md = md //? Is this used at all

	// don't remember what this one is about
	//TODO: run CheckDataDecode on debug flag
	//CheckDataDecode(G.data, md)
}

// TempData returns ptr to the central userdata editing struct var of the dscore package.
// Data must be stored in tempData to be saved (encoded/written to file) on shutdown.
// if tempData is not initialized beforehand (with InitTempData), TempData returns nil
func TempData() *globalModify {
	if tempData.initialized {
		return &tempData
	} else {
		return nil
	}
}

// InitTempData populates tempdata from globaldata
// fields populated are required to avoid data loss on toml encode
func InitTempData() {
	if !tempData.initialized {
		dat := gd.data
		tempData = globalModify{
			globalData:  &dat,
			initialized: true,
			Modified:    false,
		}
	}
}

// standardizeAlias should be applied any time a spec or component alias is set.
// It performs the following changes:
//   - removes spaces, tabs, newlines, backslash and forwardslash, and at signs
//   - converts all alphabetic to lower-case
func standardizeAlias(alias string) string {
	//TODO:(hi) Check that this cleans shit up well
	//	Switch this to trim characters from the entire string
	//	Trim all escape sequences
	return strings.ToLower(strings.Trim(alias, "\\/		'\n\r\t\""))
}

func (G *globals) forcelogG(outStr string) {
	G.GlobalMessage = append(G.GlobalMessage, outStr)
	print(outStr)
}

func (G *globals) logG(outStr string) {
	G.GlobalMessage = append(G.GlobalMessage, outStr)
}
func (G *globals) logfG(outStr string, anyfmt ...any) {
	G.GlobalMessage = append(G.GlobalMessage, fmt.Sprintf(outStr, anyfmt...))
}
