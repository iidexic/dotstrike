package pops

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"testing"
)

// ──────────────────────────────────────────────────────────────────────

func testing_job() *CopyJob {
	in := CleanPath(`d:\coding\exampleFiles\INPUT\`)
	out := CleanPath(`d:\coding\exampleFiles\OUTPUT`)
	job := cmachine.NewJob("test-job", in, out)
	// job := CopyJob{
	// 	PathIn: in, PathOut: out,
	// 	newDirs: make(map[string]bool),
	// 	BPrefs:  make(boolConfig),
	// 	SPrefs:  make(stringConfig),
	// }
	job.BPrefs[bNoRepo] = true
	job.BPrefs[bUseGlobal] = true
	job.BPrefs[bRootSubdir] = true
	job.BPrefs[bNoFiles] = false
	job.BPrefs[bNoHidden] = false

	return job

}

func testing_multiple_job() []*CopyJob {
	inpics := CleanPath(`d:\coding\exampleFiles\INPUT\pics`)
	inaudio := CleanPath(`d:\coding\exampleFiles\audio`)
	infilestruct := CleanPath(`d:\coding\exampleFiles\filestruct`)
	out := CleanPath(`d:\coding\exampleFiles\OUTPUT\Multiple`)
	out2 := CleanPath(`d:\coding\exampleFiles\OUTPUT\Multiple2`)
	jobs := make([]*CopyJob, 4)
	bprefs := make(boolConfig)
	bprefs[bNoRepo] = false
	bprefs[bUseGlobal] = false
	bprefs[bRootSubdir] = true
	bprefs[bNoFiles] = false
	bprefs[bNoHidden] = false
	jobs[0] = cmachine.NewJob("test-mult-pics", inpics, out)
	jobs[1] = cmachine.NewJob("test-mult-audio", inaudio, out)
	jobs[2] = cmachine.NewJob("test-mult-filestruct", infilestruct, out)
	jobs[3] = cmachine.NewJob("test-mult-filestruct2", infilestruct, out2)
	maps.Copy(jobs[0].BPrefs, bprefs)
	maps.Copy(jobs[1].BPrefs, bprefs)
	maps.Copy(jobs[2].BPrefs, bprefs)
	bprefs[bNoFiles] = true
	bprefs[bAllDirs] = true
	bprefs[bRootSubdir] = false
	maps.Copy(jobs[3].BPrefs, bprefs)

	return jobs
}

func readDirWalkTest(dirpath string, t *testing.T) {
	dfs := os.DirFS(dirpath)
	e := fs.WalkDir(dfs, ".", func(p string, d fs.DirEntry, e error) error {
		t.Logf("- %s (isDir: %t)", p, d.IsDir())
		return nil
	})
	if e != nil {
		t.Logf("ReadWalkDir error: %v", e)
	}
}

func testing_checkAndDeleteJobOutdir(t *testing.T, job *CopyJob, superConfirmDelete bool) {
	if job.PathOut != "" {
		cd := ReadDir(job.PathOut)
		t.Logf("CheckDir: %s", job.PathOut)
		t.Logf("DirRead: %+v", cd)

		cd.wipeDirOn = superConfirmDelete
		e := cd.deleteDir()
		if e != nil {
			t.Logf("Error deleting dir: %v", e)
		}
	}

}

func TestNonexistentPath(t *testing.T) {
	job := cmachine.NewJob("test-nonexistent-path", "D:/coding/exampleFiles/INPUT/the_secret/the_secret_2/TheMostSecretestSecretOfAll", "D:/coding/exampleFiles/OUTPUT")
	e := job.RunFS()
	if e != nil {
		t.Logf("JobRun Error (successfully didn't panic): %v", e)
	}
}

func testPathExists(t *testing.T, path string) bool {
	t.Logf("PathExists: %s", path)
	exists, e := PathExists(path)
	if exists {
		if e != nil {
			t.Logf("%s probably exists but got an error: %v", path, e)
		} else {
			t.Logf("%s exists", path)
		}
	} else {
		t.Logf("%s does not exist", path)
	}
	return exists
}

func TestRunFSdry(t *testing.T) {
	job := testing_job()
	job.BPrefs[bNoFiles] = true
	e := job.RunFS()
	if e != nil {
		t.Errorf("JobRun Error: %v", e)
	}
	t.Log(job.Detail())
	t.Log(job.DetailRunFiles())
	if len(job.OpErrors) > 0 {
		t.Log("OPERATION ERRORS:")
		for _, e := range job.OpErrors {
			t.Logf("%v", e)
		}
	}
}

func TestRunFSdirs(t *testing.T) {
	job := testing_job()
	dl, e := newDirLog(job)
	if e != nil && dl != nil {

	}
	job.BPrefs[bNoFiles] = true
	job.BPrefs[bAllDirs] = true
	e = job.RunFS()
	if e != nil {
		t.Errorf("JobRun Error: %v", e)
	}
	t.Log(job.DetailRun())

	for k, d := range job.newDirs {
		t.Logf("newDir: %s (made: %t)", k, d)
	}

	for _, e := range job.OpErrors {
		t.Errorf("%+v", e)
	}
	t.Log("DOUBLE-CHECKING OUTDIR (also deleting)")
	testing_checkAndDeleteJobOutdir(t, job, true)
}

func TestRunFSfull(t *testing.T) {
	job := testing_job()
	dl, e := newDirLog(job)
	if e != nil && dl != nil {

	}
	e = job.RunFS()
	if e != nil {
		t.Errorf("JobRun Error: %v", e)
	}
	if le := len(job.OpErrors); le > 0 {
		t.Logf("WARNING: %d OP ERRORS", le)
	}

	for _, e := range job.OpErrors {
		t.Logf("OpError %v", e)
	}
	t.Log(job.DetailRun())
	t.Log("------ POST-COPY DIR CHECK ------")
	readDirWalkTest(job.PathOut, t)
	t.Log("--- Deleting Dir Contents ---")
	e = job.wipeOutputDir()
	if e != nil {
		t.Logf("Delete-Dir Error: %v", e)
	}
	t.Log("------ DELETE ALL DIR CHECK ------")
	readDirWalkTest(job.parentPathOut, t)

}

func TestRunMultipleJobs(t *testing.T) {
	jobs := testing_multiple_job()
	for _, job := range jobs {
		e := job.RunFS()
		t.Log(job.DetailRun())
		if e != nil {
			t.Errorf("JobRun Error: %v", e)
		}
	}
	for _, job := range jobs {
		testing_checkAndDeleteJobOutdir(t, job, true)
	}
}

func TestDirFS(t *testing.T) {
	pathdir := CleanPath(`d:/coding/exampleFiles/INPUT`)

	t.Logf("DirFS at '%s'", pathdir)
	d := os.DirFS(pathdir)
	opath := "."
	f, e := d.Open(opath)
	if e != nil {
		t.Errorf("Opening '%s' Does not work.", opath)
	}
	t.Logf("Opened File:%+v", f)
	defer f.Close()
	fstat, e := f.Stat()
	if e != nil {
		t.Errorf("Stat @ '%s' failed: %v", opath, e)
	}
	t.Logf("fstat: %+v", fstat)

}

func TestDirFSWalk(t *testing.T) {
	pathdir := CleanPath(`d:/coding/exampleFiles/INPUT`)

	d := os.DirFS(pathdir)
	opath := "."
	e := fs.WalkDir(d, opath, func(p string, d DirEntry, e error) error {
		if e != nil {
			t.Logf("e @ start of walk: %s", e.Error())
		}
		info, e := d.Info()
		if e != nil {
			t.Logf("(error on Info(): %v)", e)
		}
		t.Logf("Walk: p='%s', isAbs:%t, isDir:%t, Type:%+v, Size:%v",
			p, filepath.IsAbs(p), d.IsDir(), d.Type(),
			float64(info.Size()/100000.0))
		return nil
	})
	if e != nil {
		t.Errorf("Walk Error: %v", e)
	}

}

func TestDeleteDir(t *testing.T) {
	testing_job().wipeOutputDir()
}
