package pops

import (
	"path/filepath"
	"strings"
	"testing"

	"iidexic.dotstrike/config"
)

func TestCleanPaths(t *testing.T) {
	inpath := "D:/coding/exampleFiles/INPUT"
	outpath := "D:/coding/exampleFiles/OUTPUT"
	t.Logf("WALK PATH: %s\n OUTPUT PATH: %s", inpath, outpath)
	e := filepath.WalkDir(inpath,
		func(p string, d DirEntry, e error) error {
			var rawp, winp, cleanp, relp string
			rawp = p
			cleanp = filepath.Clean(p)
			relp, rele := filepath.Rel(filepath.Clean(inpath), cleanp)
			if rele != nil {
				t.Errorf("error filepath.Rel = original:%s, cleaned:%s, relp:%s", p, cleanp, relp)
			}
			t.Logf(`
  original = '%s' | rawp = '%s'
   cleaned = '%s'
winreplace = '%s'
      relp = '%s'
---------------------`, p, rawp, cleanp, winp, relp)
			return nil
		})
	if e != nil {
		t.Errorf("%s", e.Error())
	}
}
func TestOutputPath(t *testing.T) {
	inpath := "D:/coding/exampleFiles/INPUT"
	outpath := "D:/coding/exampleFiles/OUTPUT"
	t.Logf("WALK PATH: %s\n OUTPUT PATH: %s", inpath, outpath)
	e := filepath.WalkDir(inpath,
		func(p string, d DirEntry, e error) error {
			// FUNCTION: CopyJob.stripRoot()
			t.Logf("Run filepath.Rel(%s, %s)", p, inpath)

			//t.Logf("Replace:\n%s\n%s", p, inpath)
			relp := strings.Replace(p, inpath, "", 1)
			//t.Logf("relative path = %s", relp)
			trueOutpath := filepath.Join(outpath, relp)
			t.Logf("in file: %s, out file: %s", p, trueOutpath)
			_ = trueOutpath
			return nil
		})
	if e != nil {
		t.Errorf("%s", e.Error())
	}
}

func TestCopyJob(t *testing.T) {
	cj := CopyJob{
		PathIn:  "D:/coding/exampleFiles/INPUT",
		PathOut: "D:/coding/exampleFiles/OUTPUT",
		BPrefs:  map[config.OptionKey]bool{bRootSubdir: true},
	}
	_ = cj
}

func TestPathFixes(t *testing.T) {
	inpath := "D:/coding/exampleFiles/INPUT"
	outpath := "D:/coding/exampleFiles/OUTPUT"
	//1. Condition
	inpathC, e := filepath.Abs(filepath.Clean(inpath))
	if e != nil {
		t.Logf("condition inpath:%s", e.Error())
		e = nil
	}
	outpathC, e := filepath.Abs(filepath.Clean(outpath))
	if e != nil {
		t.Logf("condition outpath:%s", e.Error())
		e = nil
	}

	//2. joinpath
	inbaseraw := filepath.Base(inpath)
	inbaseC := filepath.Base(inpathC)
	pathjoin := Joinpath(outpathC, inbaseC)
	t.Logf(`------------
     INPATH: [%s]
CONDITIONED: [%s]
    OUTPATH: [%s]
CONDITIONED: [%s]

---join---
 INPATH BASE: [%s]
INPATHC BASE: [%s]
JOINPATH: OUTPATH + INPATHC BASE
(%s) + (%s)
= %s

---Other---
JOINPATH CLEAN: %s,
FILEPATH.JOIN UNCLEAN: %s,
FILEPATH.JOIN CLEAN: %s
`, inpath, inpathC, outpath, outpathC, inbaseraw, inbaseC, outpathC, inbaseC, pathjoin, filepath.Clean(pathjoin),
		filepath.Join(outpath, inbaseraw), filepath.Join(outpathC, inbaseC))
}
