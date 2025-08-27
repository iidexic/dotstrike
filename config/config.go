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
	BoolNoFiles
	BoolCopyAllDirs
	BoolUseGlobalTarget // Spec Bools
	BoolOverrideOn
	//MaxJobCopyError  // Int (eww)
	StringGlobalTargetPath // String
	NumberOfOptions        // Count
)

type ValueType byte

const (
	TypeBool ValueType = iota
	TypeString
)

var OptionsBoolFileOp = []OptionKey{0, 1, 2, 3, 4}
var OptionsBoolSpec = []OptionKey{5, 6}
var OptionsStringGlobal = []OptionKey{7}

var ErrDecodeOptionKey = fmt.Errorf("Error finding OptionKey from decoded text")

func AllOptionIDs() []OptionKey {
	opts := make([]OptionKey, int(NumberOfOptions))
	for i := range int(NumberOfOptions) {
		opts[i] = OptionKey(i)
	}
	return opts
}

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
		Type: TypeBool, NameText: "IgnoreRepo", ForFileOp: true,
		LookupSubstrings: []string{"ignore|no", "repo|git"},
	},
	BoolIgnoreHidden: {
		Type: TypeBool, NameText: "IgnoreHidden", ForFileOp: true,
		LookupSubstrings: []string{"ignore|no", "hidden"},
	},
	BoolRootSubdir: {
		Type: TypeBool, NameText: "MakeRootSubdir", ForFileOp: true,
		LookupSubstrings: []string{"root|make", "sub"},
	},
	BoolNoFiles: {

		Type: TypeBool, NameText: "CopyNoFiles", ForFileOp: true,
		LookupSubstrings: []string{"no", "files|copy"},
	},
	BoolCopyAllDirs: {
		Type: TypeBool, NameText: "CopyAllDirs", ForFileOp: true,
		LookupSubstrings: []string{"copy|all", "all|", "dir"},
	},
	BoolUseGlobalTarget: {
		Type: TypeBool, NameText: "UseGlobalTarget", ForSpec: true,
		LookupSubstrings: []string{"use", "global", "target|tgt|"},
	},
	BoolOverrideOn: {
		Type: TypeBool, NameText: "OverrideOn", ForSpec: true,
		LookupSubstrings: []string{"override", "on|enable|"},
	},
	StringGlobalTargetPath: {
		Type: TypeString, NameText: "GlobalTargetPath", ForRun: true,
		LookupSubstrings: []string{"global|glob", "target|tgt", "path|dir"},
	},
}

func (o OptionKey) String() string { return AllOptions[o].NameText }

func (o OptionKey) IsRealOption() bool {
	if slices.Contains(AllOptionIDs(), o) {
		return true
	}
	return false
}

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

func (o OptionKey) IsBool() bool   { return AllOptions[o].Type == TypeBool }
func (o OptionKey) IsString() bool { return AllOptions[o].Type == TypeString }

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
