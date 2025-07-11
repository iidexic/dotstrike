package pops

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCleanPaths(t *testing.T) {
	inpath := "D:/coding/exampleFiles/INPUT"
	outpath := "D:/coding/exampleFiles/OUTPUT"
	t.Logf("WALK PATH: %s\n OUTPUT PATH: %s", inpath, outpath)
	e := filepath.WalkDir(inpath,
		func(p string, d DirEntry, e error) error {
			var rawp, winp, cleanp, relp string
			rawp = p
			if runtime.GOOS == "windows" {
				winp = strings.Replace(p, "/", `\`, -1)
			}
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
			if runtime.GOOS == "windows" {
				inpath = strings.Replace(inpath, "/", `\`, -1)
			}
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
		PathIn:      "D:/coding/exampleFiles/INPUT",
		PathOut:     "D:/coding/exampleFiles/OUTPUT",
		JobSettings: copyConfig{makeRootSubdir: true},
	}
	_ = cj
}

func TestPathFixes(t *testing.T) {
	inpath := "D:/coding/exampleFiles/INPUT"
	outpath := "D:/coding/exampleFiles/OUTPUT"
	//1. Condition
	inpathC, e := conditionPath(inpath)
	if e != nil {
		t.Logf("condition inpath:%s", e.Error())
		e = nil
	}
	outpathC, e := conditionPath(outpath)
	if e != nil {
		t.Logf("condition outpath:%s", e.Error())
		e = nil
	}

	//2. joinpath
	inbaseraw := filepath.Base(inpath)
	inbaseC := filepath.Base(inpathC)
	pathjoin := joinpath(outpathC, inbaseC)
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

func testCopyDir(t *testing.T, srcDir, outDir string) {
	//NOTE: Testing makeRootSubdir Code; normally just write it out. Path doesn't need to exist
	cm := GetCopierMaschine()
	cm.NewJob("test_examplefiles", srcDir, outDir)
	tcopy := cm.GetJob("test_examplefiles")
	t.Logf("CopyJob PathIn:%s, PathOut:%s", tcopy.PathIn, tcopy.PathOut)
	t.Logf("CopyJob pre-copy:\n%+v", tcopy)
	err := tcopy.Run()
	if err != nil {
		t.Errorf("COPY ERROR: %v", err)
		t.Logf("[COPYJOB: %+v]", tcopy)
	} else {
		t.Log("COPY DONE\n")
	}
	t.Logf("fstack: len = %d, # errors: %d\n--Contents:--\n", len(tcopy.fstack), len(tcopy.OpErrors))
	for i, f := range tcopy.fstack {
		if f.inSize > f.outSize {
			t.Errorf("incomplete copy error: %s", f.relpath)
		}
		t.Logf("%d. rel:%s (original=%d, new=%d)", i, f.relpath, f.inSize, f.outSize)
	}
	for i, e := range tcopy.OpErrors {
		t.Errorf("E%d: ERROR:%s", i, e.Error())
	}
	//CLEANUP
	//TODO: Clean up the dirs too
	for _, f := range tcopy.fstack {
		rmfile := filepath.Join(tcopy.PathOut, f.relpath)
		err = os.Remove(rmfile)

		if err != nil {
			t.Logf("Cleanup PathError -> `%s`", err.Error())
		}
	}
}
func TestCopyDirSimple(t *testing.T) {
	testCopyDir(t, "D:/coding/exampleFiles/INPUT", "D:/coding/exampleFiles/OUTPUT")
}

func TestCopyDirDifferentDrive(t *testing.T) {
	testCopyDir(t, "D:/coding/examplefiles", "C:/dev/.test_data/file_operations")
}
func TestCopyDirToInternal(t *testing.T) {
	testCopyDir(t, "D:/coding/exampleFiles", "D:/coding/exampleFiles/OUTPUT_INNER/")
}

func TestCopyOnlyDirSimple(t *testing.T) {
	testCopyOnlyDirs(t, "D:/coding/exampleFiles/INPUT", "D:/coding/exampleFiles/OUTPUT")
}

func testCopyOnlyDirs(t *testing.T, srcDir, outDir string) {
	cm := GetCopierMaschine()
	job1 := cm.NewJob("test_examplefiles", srcDir, outDir)
	tcopy := cm.NewJob("test_examplefiles", "", "")
	if tcopy != nil {
		t.Errorf("Failure: NewJob should return nil ptr when passed existing job name")
	} else {
		tcopy = job1 //
	}
	tcopy = cm.GetJob("test_examplefiles")
	tcopy.JobSettings.copyAllDirectories = true
	tcopy.JobSettings.noFiles = true
	t.Logf("CopyJob PathIn:%s, PathOut:%s", tcopy.PathIn, tcopy.PathOut)
	t.Logf("Job Settings:\n%+v", tcopy.JobSettings)
	err := tcopy.Run()
	if err != nil {
		t.Errorf("COPY ERROR: %v", err)
		t.Logf("[COPYJOB: %+v]", tcopy)
	} else {
		t.Log("COPY DONE\n")
	}
	t.Log("NOTE: fstack should still be written I think. this is for use as a dry run")
	t.Logf("fstack: len = %d, # errors: %d\n--Contents:--\n", len(tcopy.fstack), len(tcopy.OpErrors))
	var input_size, output_size int64 = 0, 0
	for _, f := range tcopy.fstack {
		input_size += f.inSize
		output_size += f.outSize
	}
	if output_size > 0 {
		t.Error("FAIL - OUTPUT SIZE SHOULD BE 0")
	}
	t.Logf("Size sum: In = %d, Out = %d", input_size, output_size)
	t.Logf("Dirs:\n")
	for path, wroteDir := range tcopy.newDirs {
		t.Logf("(%s) added = %t", path, wroteDir)

	}
	for i, e := range tcopy.OpErrors {
		t.Errorf("E%d: ERROR:%s", i, e.Error())
	}
	//CLEANUP
	//currently get errors because can't delete dir while it contains another dir
	for kpath := range tcopy.newDirs {
		rmdir := filepath.Join(tcopy.PathOut, kpath) //kpath is relative dir path
		err = os.Remove(rmdir)

		if err != nil {
			t.Logf("Cleanup PathError -> `%s`", err.Error())
		}
	}
}
