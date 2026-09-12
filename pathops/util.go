package pops

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func endswith(s, suffix string) bool {
	ls := len(s)
	lsuf := len(suffix)
	if ls >= lsuf {
		return s[ls-lsuf:] == suffix
	}
	return false
}

func startswith(s, prefix string) bool {
	if lp := len(prefix); len(s) >= lp {
		return s[:lp] == prefix
	}
	return false
}

// "D:/coding/exampleFiles/OUTPUT" -> ["D:", "coding", "exampleFiles", "OUTPUT"]
func SplitAbsPath(path string) []string {
	path = CleanPath(path)
	return strings.Split(path, `\`)
}

func stripRoot(root, path string) (string, error) {
	relpath, e := filepath.Rel(root, path)
	if e != nil {
		return relpath, e
	}
	return relpath, nil
}

func PathsMatch(p1, p2 string) bool {
	p1, p2 = CleanPath(p1), CleanPath(p2)
	return p1 == p2
}

// takes in time.Time and returns date string (yy-mm-dd) and time string (hh:mm)
func DateTimeDetail(t time.Time) (string, string) {
	dy, dm, dd := t.Date()
	th, tm, _ := t.Clock()
	date := fmt.Sprintf("%02d-%02d-%02d", dy, dm, dd)
	time := fmt.Sprintf("%02d:%02d", th, tm)
	return date, time
}

// addRelpath will add the root local portion of 'from' onto 'addto'
// it is akin to filepath.Rel()
func addRelpath(addto, root, from string) (string, error) {
	addto, root, from = CleanPath(addto), CleanPath(root), CleanPath(from)
	if CleanPath(root) == CleanPath(from) {
		return addto, nil
	}
	rel, e := filepath.Rel(root, from)
	if e != nil {
		if rel == "" {
			return "", fmt.Errorf("Rel failed (Rel(%s,%s)) [%w]", root, from, e)
		}
	}
	// basically the same as `if root==from` above. Can probably remove
	if rel == "." {
		return addto, e
	}
	return Joinpath(addto, rel), e
}
