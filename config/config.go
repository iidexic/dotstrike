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

// option spec contains required information for each option.
// includes name, type, use/purpose, and lookup string slices
// LookupSubstrings uses `|` to indicate separate values that can be used in the same place
type option struct {
	Type             ValueType
	LookupSubstrings []string
	NameText         string
	ForFileOp        bool
	ForRun           bool
	ForSpec          bool
}

var AllOptions = map[OptionKey]option{
	BoolIgnoreRepo: {
		Type: Tbool, NameText: "IgnoreRepo", ForFileOp: true,
		LookupSubstrings: []string{"ignore|no", "repo|git"},
	},
	BoolIgnoreHidden: {
		Type: Tbool, NameText: "IgnoreHidden", ForFileOp: true,
		LookupSubstrings: []string{"ignore|no", "hidden"},
	},
	BoolRootSubdir: {
		Type: Tbool, NameText: "MakeRootSubdir", ForFileOp: true,
		LookupSubstrings: []string{"root|make", "sub"},
	},
	BoolSourceSubdirs: {
		Type: Tbool, NameText: "SourceSubdirs", ForFileOp: true,
		LookupSubstrings: []string{"source|src", "sub|dirs"},
	},
	BoolNoFiles: {

		Type: Tbool, NameText: "CopyNoFiles", ForFileOp: true,
		LookupSubstrings: []string{"no", "files|copy"},
	},
	BoolCopyAllDirs: {
		Type: Tbool, NameText: "CopyAllDirs", ForFileOp: true,
		LookupSubstrings: []string{"copy|all", "all|", "dir"},
	},
	BoolUseGlobalTarget: {
		Type: Tbool, NameText: "UseGlobalTarget", ForSpec: true,
		LookupSubstrings: []string{"use", "global", "target|tgt|"},
	},
	BoolKillGlobalTarget: {
		Type: Tbool, NameText: "DisableGlobalTarget", ForSpec: true,
		LookupSubstrings: []string{"kill|disable", "global", "target|tgt|"},
	},
	BoolOverrideOn: {
		Type: Tbool, NameText: "OverrideOn", ForSpec: true,
		LookupSubstrings: []string{"override", "on|enable|"},
	},
	StringGlobalTargetPath: {
		Type: Tstring, NameText: "GlobalTargetPath", ForRun: true,
		LookupSubstrings: []string{"set", "glob|glb", "target|tgt|dest", "path|dir"},
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
