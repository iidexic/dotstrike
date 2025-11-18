package pops

import (
	"strings"
)

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

type matchPath interface {
	matches(string) bool
	appliesIf(bool) bool
}

// IgnoreSet stores & processes ignore data for a CopyJob
type IgnoreSet struct {
	Patterns []matchPath
}

// type prefixptn struct { ptn string; matchDir, matchFile bool; inRootDir bool }
// func (pp prefixptn) matches(s string) bool {/* did not finish this */ }
// func (pp prefixptn) appliesIf(isDir bool) bool { return (isDir && pp.matchDir || !isDir && pp.matchFile) }

// subptn is a single ignore string pattern
// requires subptn.pattern string and min 1 of matchDir, matchFile
type subptn struct {
	ptn                 string
	matchDir, matchFile bool
	//anyL, anyR bool
	//psize      byte
}

func (ip subptn) appliesIf(isDir bool) bool { return (isDir && ip.matchDir || !isDir && ip.matchFile) }

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
		if pat.appliesIf(isDir) {
			if pat.matches(path) {
				return true
			}
		}
	}
	return false
}

// type classifyPattern struct {
// 	ptn, hardpath string
// 	splitpath     []string
// }

// Adds a subptn to the IgnoreSet. A subptn only checks whether the pattern string exists as a substring anywhere within a given path string
func (I *IgnoreSet) AddSubpattern(ptn string, matchDir, matchFile bool) {
	I.Patterns = append(I.Patterns, subptn{ptn: ptn, matchDir: matchDir, matchFile: matchFile})

}
