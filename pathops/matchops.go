package pops

import (
	"strings"
)

// ── Ignore Functionality ────────────────────────────────────────────
// TODO: Make + flesh out a string pattern matching package. Useful for all 3 main packages

// matchString struct reports whether an input string matches itself

// ----------- thinking thru tokenizing ------------
// type ptok rune
// var (
//
//	localdir ptok='.'
//	parentdir ptok ='.'
//	usep ptok = '\\'
//	wsep ptok = '/'
//	home ptok = '~'
//	wild ptok = '*'
//
// )
// -------------------------------------------------
type matchStyle int

const (
	ptnInVal matchStyle = iota
	valInPtn
	exactlyValNoCase
	exactlyVal
	isValPathName
	isValParentName
	inValPathName
	inValParentName
)

type matchString interface {
	matches(string) bool
}

// IgnoreSet stores & processes ignore data for a CopyJob
type IgnoreSet struct {
	Patterns []subptn
}

// subptn is a single ignore string pattern
// requires subptn.pattern string and min 1 of matchDir, matchFile
type subptn struct {
	ptn                 string
	tokens              []rune
	matchDir, matchFile bool
	//anyL, anyR bool
	//psize      byte
}

// matches checks string against the valid iptn
func (ip subptn) matches(s string) bool {
	//TODO: Requires Improvement. Account for Path separators

	if strings.Contains(s, ip.ptn) {
		return true
	}
	return false
}

func (I *IgnoreSet) isIgnored(path string, isDir bool) bool {

	for _, pat := range I.Patterns {
		if (!isDir && pat.matchFile) || (isDir && pat.matchDir) {
			if pat.matches(path) {
				return true
			}
		}
	}
	return false
}
