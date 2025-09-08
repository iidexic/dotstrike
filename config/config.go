package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

type OptionKey int

const (
	NotAnOption OptionKey = iota - 1
	BoolIgnoreRepo
	BoolIgnoreHidden
	BoolRootSubdir
	BoolSourceSubdirs
	BoolNoFiles
	BoolCopyAllDirs
	BoolUseGlobalTarget // Spec Bools
	BoolKillGlobalTarget
	BoolOverrideOn
	StringGlobalTargetPath // String
	NumberOfOptions        // Count
	//MaxJobCopyError  // Int - requires implementation
)

type ValueType byte

const (
	_ ValueType = iota
	Tbool
	Tstring
	TstringSlice
)

/*
WARNING: As is, adding new options requires
1. Adding to OptionKey Enum
2. AllOptions; with an option containing all required option info
*/

// var OptionsBoolFileOp = []OptionKey{0, 1, 2, 3, 4}
// var OptionsBoolSpec = []OptionKey{5, 6}
// var OptionsStringGlobal = []OptionKey{7}

var ErrDecodeOptionKey = fmt.Errorf("Error finding OptionKey from decoded text")

func AllOptionIDs() []OptionKey {
	opts := make([]OptionKey, int(NumberOfOptions))
	for i := range int(NumberOfOptions) {
		opts[i] = OptionKey(i)
	}
	return opts
}

//TODO:(HIGHEST) REPLACE LOOKUPS WITH DEFINED FLAGS WHEREVER THEY ARE HAPPENING

// option spec contains required information for each option.
// includes name, type, use/purpose, and lookup string slices
// LookupSubstrings uses `|` to indicate separate values that can be used in the same place
type option struct {
	Type             ValueType
	LookupSubstrings []string
	LookupExacts     []string
	NameText         string
	fName, fshort    string
	runUsage         string
	ForFileOp        bool
	ForRun           bool
	ForSpec          bool
}

var AllOptions = map[OptionKey]option{
	BoolIgnoreRepo: {
		Type: Tbool, NameText: "IgnoreRepo", fName: "ignore-repo",
		runUsage:  "Disables copy of the .git repo by adding to global ignores.",
		ForFileOp: true, LookupSubstrings: []string{"ignore|no", "repo|git"},
		LookupExacts: []string{"nore"},
	},
	BoolIgnoreHidden: {
		Type: Tbool, NameText: "IgnoreHidden", fName: "ignore-hidden",
		runUsage:  `Add hidden paths to global ignores; Disables copy of dir/file names that begin with '_' or '.'`,
		ForFileOp: true, LookupSubstrings: []string{"ignore|no", "hidden"}, LookupExacts: []string{"nohi"},
	},
	BoolRootSubdir: {
		Type: Tbool, NameText: "MakeRootSubdir", fName: "make-subdir",
		runUsage: `Makes a new dir in target folder to copy a spec into.
Dir is named with spec's alias if possible, else numbers will be added`,
		ForFileOp: true, LookupSubstrings: []string{"root|make", "root|sub", "dir"},
		LookupExacts: []string{"mrsd"},
	},
	BoolSourceSubdirs: {
		Type: Tbool, NameText: "SourceSubdirs", fName: "separate-sources",
		runUsage:  "Copies each source into a separate subdir; name is source's alias or source path's dir name.",
		ForFileOp: true, LookupSubstrings: []string{"source|src", "sub|dirs"},
		LookupExacts: []string{"ssep"},
	},
	BoolNoFiles: {
		Type: Tbool, NameText: "CopyNoFiles", fName: "no-files", fshort: "n",
		runUsage:  "Disable filecopy for run. Use for dry runs, or with --all-dir to copy only the directory structure",
		ForFileOp: true, LookupSubstrings: []string{"no", "files|copy"}, LookupExacts: []string{"dryrun", "dry"},
	},
	BoolCopyAllDirs: {
		Type: Tbool, NameText: "CopyAllDirs", fName: "all-dirs", fshort: "d",
		runUsage: `Copy all Source subdirectories, including empty subdirectories. 
Use with --no-files to only copy the directories themselves.`,
		ForFileOp: true, LookupSubstrings: []string{"copy|all", "all|", "dir"}, LookupExacts: []string{"alldirs", "aldr"},
	},
	BoolUseGlobalTarget: {
		Type: Tbool, NameText: "UseGlobalTarget", fName: "use-global-target",
		runUsage: `Use "--GlobalTarget" to enable write to Global Target for all specs in run.
Use --GlobalTarget="off" to forcibly disable write to GlobalTarget for full run, including specs that exclusively target the Global Target.`,
		ForSpec: true, LookupSubstrings: []string{"use|all|force", "global|glb|gtg", "target|tgt|"}, LookupExacts: []string{"usegt", "allgt", "agtg"},
	},
	BoolKillGlobalTarget: {
		Type: Tbool, NameText: "DisableGlobalTarget", fName: "kill-global-target",
		runUsage: `--GlobalTarget=off`,
		ForSpec:  true, LookupSubstrings: []string{"kill|disable|no", "glob|glb|gtg", "target|tgt|"},
		LookupExacts: []string{"nogt", "gtoff", "gtgoff"},
	},
	BoolOverrideOn: {
		Type: Tbool, NameText: "OverrideOn", fName: "use-override",
		runUsage: ``,
		ForSpec:  true, LookupSubstrings: []string{"|use", "override", "on|enable|"}, LookupExacts: []string{"over", "esor"},
	},
	StringGlobalTargetPath: {
		Type: Tstring, NameText: "GlobalTargetPath", fName: "set-global-target-path",
		runUsage: ``, ForRun: true,
		LookupSubstrings: []string{"set", "glob|glb", "target|tgt|dest", "path|dir|"},
		LookupExacts:     []string{"gtpath", "globaltarget", "gtpath"},
	},
}

// TODO: Figure out if really want to make this a val, error return or just make 100% SURE tests get run

func (o OptionKey) String() string {
	_, ok := AllOptions[o]
	if ok {
		return AllOptions[o].NameText
	}
	return "FAILURE_OPTION_NOT_FOUND_IN_ALLOPTIONS"
}
func (o OptionKey) RunUsage() string   { return AllOptions[o].runUsage }
func (o OptionKey) NameFlag() string   { return AllOptions[o].fName }
func (o OptionKey) NameFshort() string { return AllOptions[o].fshort }
func (o OptionKey) IsRealOption() bool { return slices.Contains(AllOptionIDs(), o) }

func (o OptionKey) MarshalTOML() ([]byte, error) {
	return toml.Marshal(o.String())
}
func (o *OptionKey) UnmarshalTOML(data []byte) error {
	var text string
	if e := toml.Unmarshal(data, &text); e != nil {
		return e
	}
	outopt := LookupOption(text)
	if outopt != NotAnOption {
		o = &outopt
		return nil
	}
	return ErrDecodeOptionKey
}

func (o OptionKey) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}
func (o *OptionKey) UnmarshalText(data []byte) error {
	optkey := LookupOption(string(data))
	if optkey.IsRealOption() {
		o = &optkey
	}
	return ErrDecodeOptionKey
}
func OptFrom(optionName string) OptionKey {
	for k, v := range AllOptions {
		if optionName == v.NameText {
			return k
		}
	}
	return NotAnOption
}

func (o OptionKey) IsBool() bool   { return AllOptions[o].Type == Tbool }
func (o OptionKey) IsString() bool { return AllOptions[o].Type == Tstring }

// TODO: (mid-fix) maybe count up the number of matches to ensure only 1? Would require other changes
func LookupOption(input string) OptionKey {
	input = strings.TrimSpace(strings.ToLower(input))
	for id, opt := range AllOptions {
		match := true
		for _, substr := range opt.LookupSubstrings {
			match = match && lookupSubstringMatch(input, substr)
		}
		if match {
			return id
		}
	}
	return NotAnOption
}

func LookupOptionExact(input string) (OptionKey, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	for id, opt := range AllOptions {
		if slices.Contains(opt.LookupExacts, input) {
			return id, nil
		}
	}
	return NotAnOption, nil
}

// func div(num, divis int) float32 {
// 	if divis != 0 {
// 		return float32(num) / float32(divis)
// 	}
// 	return -0
// }

// Returns 1 key for each string. If string does not match, returns NotAnOption
func GetOptionKeys(searches []string) []OptionKey {
	found := make([]OptionKey, len(searches))
	for i, s := range searches {
		found[i] = LookupOption(s)
	}
	return found
}
func SimplestSearchString(opt OptionKey) string {
	return strings.Join(firstsubs(AllOptions[opt].LookupSubstrings), "")
}
func firstsubs(lookupsubs []string) []string {
	out := make([]string, 0, len(lookupsubs)*2)
	for _, sub := range lookupsubs {
		out = append(out, firstsub(sub))
	}
	return out
}

func firstsub(sub string) string {
	if strings.Contains(sub, "|") {
		return sub[:strings.Index(sub, "|")]
	}
	return sub
}

func lookupSubstringMatch(input string, sub string) bool {
	// break up substring if has an or
	//sublist := make([]string, 0, 2)
	for s := range strings.SplitSeq(sub, "|") {
		if strings.Contains(input, s) {
			return true
		}
	}
	return false
}
