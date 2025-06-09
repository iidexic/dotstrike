package dscore

import (
	"fmt"
	"strings"
)

type tomlVar interface {
	string | int | float64
}

type tomlTable[T tomlVar] []T

type tconverter struct {
	tlines []string
}

func fromTomlString(tomlstr string) {
	toml := tomlCleanLines(tomlstr)
	tc := newTomlParser(toml)
	_ = tc
	//print("splitting:", tomlstr)
}

// Split the toml string to lines, remove comments, trim
func tomlCleanLines(tomlstr string) []string {
	lntoml := strings.Split(tomlstr, "\n")
	print("|-TomlSplit-|\n")
	for i, ln := range lntoml {
		fmt.Printf("[%d] %s\n", i, ln)
		pos := strings.Index(ln, "#")
		if pos >= 0 {
			lntoml[i] = strings.Trim(ln[:pos], " ")
		}
	}
	return lntoml
}

func newTomlParser(dataLines []string) *tconverter {
	t := tconverter{
		tlines: dataLines,
	}
	return &t
}

/*
Data from toml:
[global]
storagePath = '~/.config/dotstrike/store'
[global.prefs]
keep.repo = true
keep.hidden = true
storedata.sourcedir = true
storedata.central = true

[global.data]
[[cfgs]]
alias = 'wezterm'
sources = ['~/.config/wezterm']
targets = ['@GLOBAL']

[[cfgs]]
alias = 'nvim'
sources = ['~/AppData/Local/nvim']
targets = ['@GLOBAL']
*/
