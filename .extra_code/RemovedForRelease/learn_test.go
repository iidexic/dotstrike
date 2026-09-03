package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pops "iidexic.dotstrike/pathops"
)

// ╭─────────────────────────────────────────────────────────╮
// │        Path tests to check general go behaviour         │
// ╰─────────────────────────────────────────────────────────╯

var splitList = filepath.SplitList
var testpathsh = []string{"D:/coding/exampleFiles/INPUT/file_format/bad-gif.gif", "D:/coding/exampleFiles/OUTPUT", "./thing", ".../a//b\\\\\\c"}

func splitPathNotDumb(path string) []string {
	if strings.Contains(path, `\`) {
		return strings.Split(path, `\`)
	}
	if strings.Contains(path, `/`) {
		return strings.Split(path, `/`)
	}
	return []string{path}
}

func CheckDir(t *testing.T, path string) {
	t.Logf("CheckDir: %s", path)
	d, e := os.Stat(path)
	if e != nil {
		t.Logf("Error getting stat: %v", e)
	}
	if d.IsDir() {
		t.Logf("%s exists as dir", path)
	} else {
		t.Logf("%s exists as file", path)
	}
	e = fs.WalkDir(os.DirFS(path), ".", func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			t.Logf("WalkDir error: %v", e)
		}
		if d.IsDir() {
			t.Logf("WalkDir: %s is a dir", p)
		}
		t.Logf("path: %s", p)
		return nil
	})

}

func TestPathSplits(t *testing.T) {
	for i, p := range testpathsh {
		sl := splitList(p)
		snd := splitPathNotDumb(p)
		t.Logf("[%d]SplitPathNotDumb(len %d): %v", i, len(snd), snd)
		t.Logf("SplitList(len %d): %v", len(sl), sl)
	}
}

func TestCleanPaths(t *testing.T) {
	ps := []string{"./gamer", "../gamer", "...\\gamer", "~/gamer"}
	for _, p := range ps {
		t.Logf("CleanPath(%s) = %s", p, pops.CleanPath(p))
	}
}

func testPathSplitting(t *testing.T, path string) {
	t.Logf("path:%s", path)
	path = pops.CleanPath(path)
	t.Logf("clean:%s", path)
	countBS := strings.Count(path, `\`)
	countFS := strings.Count(path, "/")
	t.Logf("countBS:%d, countFS:%d", countBS, countFS)
	plist := make([]string, max(countBS, countFS)+1)

	for i := len(plist) - 1; i >= 0; i-- {
		subp, basep := filepath.Split(path)
		plist[i] = basep
		t.Logf("plist[%d]: %s (path: %s, subp: %s)", i, plist[i], path, subp)
		//path = filepath.Clean(subp)
		if strings.HasSuffix(subp, `\`) || strings.HasSuffix(subp, `/`) {
			path = subp[:len(subp)-1]
		} else {
			path = subp
		}
		if i == 1 {
			plist[0] = subp
		}
	}
}

func TestPathSplitLogic(t *testing.T) {
	testPathSplitting(t, "D:/coding/exampleFiles/OUTPUT")
	testPathSplitting(t, "./thing")
	testPathSplitting(t, ".../a//b\\\\\\c")
	testPathSplitting(t, "D:/GamerZone/")
	testPathSplitting(t, "D/GamerZone")
}

func TestPathExplore(t *testing.T) {
	pfile := "D:/coding/exampleFiles/INPUT/file_format/bad-gif.gif"
	split := filepath.SplitList(pfile)
	clean := filepath.Clean(pfile)
	t.Logf("clean:%s\nsplit:%v", clean, split)
	localize, err := filepath.Localize(filepath.Clean(pfile))
	t.Logf("Localized = %s (err %v)", localize, err)
}
