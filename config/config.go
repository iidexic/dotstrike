package config

import (
	"slices"
	"strings"
)

type OptionKey int

const (
	NotAnOption OptionKey = iota - 1
	BoolIgnoreGit
	BoolIgnoreHidden
	BoolRootSubdir
	BoolNoFiles
	BoolCopyAllDirs
	BoolUseGlobalTarget // Spec Bools
	BoolOverrideOn
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
	BoolIgnoreGit: {
		Type:             TypeBool,
		LookupSubstrings: []string{"ignore|no", "repo|git"},
		NameText:         "IgnoreRepo", ForFileOp: true,
	},
	BoolIgnoreHidden: {
		Type:             TypeBool,
		LookupSubstrings: []string{"ignore|no", "hidden"},
		NameText:         "IgnoreHidden", ForFileOp: true,
	},
	BoolRootSubdir: {
		Type:             TypeBool,
		LookupSubstrings: []string{"root|make", "subdir"},
		NameText:         "MakeRootSubdir", ForFileOp: true,
	},
	BoolCopyAllDirs: {
		Type: TypeBool, NameText: "CopyAllDirs", ForFileOp: true,
		LookupSubstrings: []string{"all", "dir|dirs"},
	},
	BoolOverrideOn: {
		Type: TypeBool, NameText: "OverrideOn", ForSpec: true,
		LookupSubstrings: []string{"override", "on|enable"},
	},
	BoolUseGlobalTarget: {
		Type: TypeBool, NameText: "UseGlobalTarget", ForRun: true,
		LookupSubstrings: []string{"use", "global", "target"},
	},
}

func (o OptionKey) Text() string { return AllOptions[o].NameText }

func (o OptionKey) IsRealOption() bool {
	if slices.Contains(AllOptionIDs(), o) {
		return true
	}
	return false
}

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
