package dscore

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

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
	Selected         int    `toml:"SelectedSpec"`
	GlobalTargetPath string `toml:"targetpath"`
	Prefs            prefs  `toml:"prefs"`
	Specs            []Spec `toml:"specs"` // possibly needs omitempty
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

// !TODO:(hi-refactor) Change to use maps for prefs. Either make prefs a map or make prefs contain maps
// - This will majorly simplify working with config options
type prefs struct {
	bools map[ConfigOption]bool
	//TODO: symlink handling + symlink preference
}

// ╭─────────────────────────────────────────────────────────╮
// │                     CONFIG OPTIONS                      │
// ╰─────────────────────────────────────────────────────────╯

// uses ConfigOption index
type ConfigOption int

// TODO:(mid-feat) better system (when do patterns/regex)

// ── ConfigOptions ───────────────────────────────────────────────────

const (
	NotAnOption ConfigOption = iota - 1
	OptBKeepRepo
	OptBKeepHidden
	OptBUseGlobalTgt
	OptBCopyFiles
	OptBCopyAllDirs
	OptSGlobalTargetPath
)

// BoolOptions is a list containing all bool-value ConfigOptions
// requires all val == index
var BoolOptions = []ConfigOption{0, 1, 2, 3, 4}

// StringOptions is a list containing all string-value ConfigOptions
// requires all val == index + len(BoolOptions)
var StringOptions = []ConfigOption{5}

// var AllOptions = append(BoolOptions,StringOptions...)

// optionCount provides the total length of the ConfigOption enum
var optionCount = len(BoolOptions) + len(StringOptions)

var PrefIdentifiers = [][]string{
	{"keeprepo", "keep-repo", "keep_repo", "repo"},
	{"keephidden", "keep-hidden", "keep_hidden", "hidden"},
	{
		"useglobaltarget", "use-global-target", "use_global_target",
		"useglobaltgt", "use-globaltarget",
		"globaltarget", "use-global",
	},
	{"copyfiles", "copy-files", "copy_files", "docopy"},
	{
		"alldirs", "all-dirs", "all_dirs",
		"alldir", "all-dir", "all_dir",
		"copyalldirs", "copy-all-dirs", "copy-alldir",
	},
	{"globaltargetpath", "targetpath", "global_target_path", "global-target-path"},
}

func (c ConfigOption) Text() string {
	switch c {
	case OptBKeepHidden:
		return "KeepHidden"
	case OptBKeepRepo:
		return "KeepRepo"
	case OptBUseGlobalTgt:
		return "UseGlobalTarget"
	case OptBCopyFiles:
		return "CopyFiles"
	case OptBCopyAllDirs:
		return "CopyAllDirs"
	case OptSGlobalTargetPath:
		return "GlobalTargetPath"
	}
	return ("NotAnOption")
}

// OptionID returns the ConfigOption that optName corresponds to within PrefIdentifiers.
//
// optName is run through quickclean() before checking available values.
func OptionID(optName string) ConfigOption {
	//OptionID relies on having ConfigOption Enum values match PrefIdentifiers index
	for i := range PrefIdentifiers {
		if slices.Contains(PrefIdentifiers[i], quickclean(optName)) {
			return ConfigOption(i)
		}

	}
	return NotAnOption
}

func IsBoolOption(opt ConfigOption) bool { return slices.Contains(BoolOptions, opt) }

func IsStringOption(opt ConfigOption) bool { return slices.Contains(StringOptions, opt) }

// ──────────────────────────────────────────────────────────────────────

func (G *globals) Detail() string {
	lines := make([]string, 1, 32)
	lines[0] = fmt.Sprintf(`GLOBAL USER DATA:
=================
Config Path: '%s'
Selected Spec(index): %d

Globals Log (instance):
`, G.dsconfigPath, G.data.Selected)
	lines = append(lines, G.GlobalMessage...)

	lines = append(lines, G.data.Prefs.Detail())
	for i := range G.data.Specs {
		lines = append(lines, G.data.Specs[i].Detail())
	}

	/* need
	0. dscdsconfigPath
	1. GlobalMessage
	2.
	*/
	return strings.Join(lines, "\n")
}

func (p prefs) Detail() string {
	out := ""
	for k, v := range p.bools {
		out = fmt.Sprintf("%s\n%s:%t", out, k.Text(), v)
	}
	return out
}

func (p prefs) equal(p2 prefs) bool {
	return maps.Equal(p.bools, p2.bools)
}

// TempGlob exists to store new global data temporarily during runtime
// this will then be checked/merged with GD, and written to globals file
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
	return strings.ToLower(strings.Trim(alias, "\\/		\n@"))
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
	file.Truncate(0)
	encode := toml.NewEncoder(file)
	e = encode.Encode(gm.globalData)
	if e != nil {
		return e
	} else {
		return nil
	}
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
