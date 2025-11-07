package pops

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"iidexic.dotstrike/config"
	"iidexic.dotstrike/uout"
)

type fileRecord struct {
	files  map[string]*filedata
	outdir string
}

// TODO:(low) better system for bools at least (enum?)
type filedata struct {
	//job                            *CopyJob
	preExisting, copied, isDir     bool
	ignored                        bool
	sizeOrigin, sizePre, sizeFinal int64
	dataWritten                    int64
	modOrigin, modFinal            time.Time
}

func (R fileRecord) String() string {
	out := uout.NewOut("Files Seen:")
	for path, data := range R.files {
		if !data.isDir {
			out.NV(path, data)
		}
	}
	return out.String()
}

func (F filedata) String() string {
	if F.isDir {
		return ""
	}
	out := uout.NewOut("")
	switch {
	case F.preExisting && F.copied: // copied over file
		out.A("Copied over existing")
	case F.ignored && F.copied:
		out.A("Well it got ignored but that didn't work because it still got copied")
	case F.ignored:
		out.A("Ignored")
	case F.preExisting && F.isDir: // not copied: file exists
		out.A("Existing Dir")
	case F.copied: // copied, new
		out.A("Copied, new")
	}
	if !F.isDir {
		out.A(fmt.Sprintf(" (%.2f%% copied)", F.percentSize()*100))
	}
	return out.String()
}

func (F filedata) percentSize() float64 {
	switch {
	case F.sizeOrigin == F.sizeFinal && F.sizeOrigin > 0:
		return 1.0
	case F.sizeOrigin == 0 && F.sizeFinal > 0:
		return -1
	case F.sizeOrigin == 0:
		return 0
	}
	return float64(F.sizeFinal) / float64(F.sizeOrigin)

}

// newRecord creates a record for the given path.
// It runs os.Stat on outpath to get the status before copy run.
func (R *fileRecord) newRecord(outpath, relpath string, isDir bool) (*filedata, error) {
	R.files[relpath] = &filedata{isDir: isDir}
	fd := R.files[relpath]
	info, err := os.Stat(outpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fd.preExisting = false
			return fd, nil
		}
		return fd, err
	}
	fd.preExisting = true
	if !isDir {
		fd.sizePre = info.Size()
	}
	return fd, nil
}

func (F *filedata) setOrigin(info fs.FileInfo) {
	F.sizeOrigin = info.Size()
	F.modOrigin = info.ModTime()

}

func (F *filedata) setNew(info fs.FileInfo, written int64) {
	F.sizeFinal = info.Size()
	F.modFinal = info.ModTime()
	F.copied = true
	F.dataWritten = written
}

// CopyJob prepares and executes the copy of all contents of PathIn to PathOut
type CopyJob struct {
	PathIn, PathOut string     // Root of copy source and destination (*or destination parent)
	parentPathOut   string     // unused. populated on run if JobSettings.makeRootSubdir = true
	fstack          []filecopy // record of files copied
	record          fileRecord
	jobRan, abort   bool
	newDirs         map[string]bool //
	ignore          IgnoreSet
	OpErrors        []fs.PathError
	BPrefs          boolConfig
	SPrefs          stringConfig //* Currently Unused
}

func (J CopyJob) String() string {
	return J.Detail()
}

func (J *CopyJob) DetailLine() string {
	d := fmt.Sprintf("in:'%s' | out: %s ", J.PathIn, J.PathOut)
	if len(J.ignore.Patterns) > 0 {
		d += fmt.Sprintf("| #ignores:%d", len(J.ignore.Patterns))
	}
	if J.jobRan {
		d += fmt.Sprintf("| ran (%d file, %d newdir, %d errors)", len(J.record.files), J.DirsMade(), len(J.OpErrors))
	}
	return d
}

func (J *CopyJob) Pref(opt config.OptionKey) bool {
	b, ok := J.BPrefs[opt]
	return ok && b
}

func (J *CopyJob) detailOutpath() string {
	if J.parentPathOut != "" {
		base := Base(J.PathIn)
		return fmt.Sprintf("'%s' in '%s\\'", J.PathOut, base)
	}
	return fmt.Sprintf("'%s'", J.PathOut)
}

func (J *CopyJob) Detail() string {
	out := uout.NewOut("Job: ")
	out.IndR()
	out.F("in: '%s' out: %s", J.PathIn, J.detailOutpath())
	out.IfV(J.jobRan, "Job Ran: ", "Job did not run")
	if J.jobRan {
		nf, ad := J.Pref(bNoFiles), J.Pref(bAllDirs)
		switch {
		case nf && ad:
			out.A("(dir structure only)")
		case nf:
			out.A("(dry run)")
		case ad:
			out.A("(all dirs)")
		}
		out.IndR()
		out.F("%d files copied (%.2f%% dir data),  %d dirs seen, ~%d dirs made",
			len(J.record.files), J.CopyPercent()*100, len(J.newDirs), J.DirsMade())
		out.IndL()
	}
	out.V("Job Preferences:")
	out.IndR()
	out.LnSplit(config.DetailFlat(J.BPrefs))
	return out.String()
}

func (J CopyJob) DetailRun() string {
	out := uout.NewOut("Job: ")
	bin := Base(J.PathIn)
	var bout string
	if J.parentPathOut != "" {
		bout = Base(J.parentPathOut)
	} else {
		bout = Base(J.PathOut)
	}
	out.F("From '%s'  to '%s'", bin, bout)
	out.Ln(J.DetailRunFiles())
	out.Ln(J.DetailRunDirs())
	return out.String()
}

func (J CopyJob) DetailRunFiles() string {
	return J.record.String()
}

func (J CopyJob) DetailRunDirs() string {
	out := uout.NewOut("Directories")
	out.IndR()
	for k, v := range J.newDirs {
		out.F("'%s': made=%t", k, v)
	}

	return out.String()
}

func (J CopyJob) CopyPercent() float64 {
	var read, copied int64 = 0, 0
	for _, f := range J.record.files {
		read += f.sizeOrigin
		copied += f.sizeFinal
	}
	switch copied {
	case 0:
		return 0.0
	case read:
		return 1.0
	default:
		return float64(copied) / float64(read)
	}
}

func (J *CopyJob) logError(abspath, opname string, e error) {
	// should I check or wrap it anyway?
	pe, ok := any(e).(PathError)
	if ok {
		J.OpErrors = append(J.OpErrors, pe)
	} else {
		J.OpErrors = append(J.OpErrors, fs.PathError{Path: abspath, Op: opname, Err: e})
	}
}

// checkAndLogError checks the error, and logs non-nil errors to CopyJob.logError.
// returns true if error!=nil, else false
func (J *CopyJob) checkAndLogError(abspath, opname string, e error) bool {
	if e != nil {
		J.logError(abspath, opname, e)
		return true
	}
	return false
}

// logDir adds directories to j.newDirs if they are not already present
// NOTE: Walk sends relative paths to logDir (J.newDirs keys will be relative)
func (J *CopyJob) logDir(dir string, copied bool) {
	var exists bool
	for keydir := range J.newDirs {
		exists = (exists || keydir == dir)
	}
	if !exists {
		J.newDirs[dir] = copied
	}
}
func (J *CopyJob) DirsMade() int {
	n := 0
	for _, v := range J.newDirs {
		if v {
			n++
		}
	}
	return n
}

func (J *CopyJob) configCheck(opt config.OptionKey) bool {
	if opt.IsBool() {
		v, ok := J.BPrefs[opt]
		return v && ok
	}
	if opt.IsString() {
		v, ok := J.SPrefs[opt]
		return len(v) > 0 && ok
	}

	return false
}
func (J *CopyJob) wipeOutputDir() error {
	wpath := J.PathOut
	if J.configCheck(bRootSubdir) {
		e := os.RemoveAll(wpath)
		return e
	}
	dir := os.DirFS(wpath)
	gobby, eeror := fs.Glob(dir, `.\*`)
	if eeror != nil {
		return eeror
	}
	var eout error
	for _, g := range gobby {

		e := os.RemoveAll(g)
		if e != nil {
			if eout == nil {
				eout = fmt.Errorf("RemoveAll Errors: %w", e)
			} else {
				eout = fmt.Errorf("%w, %w", eout, e)
			}
		}
	}
	return eout
}

type checkDir struct {
	dirpath   string
	record    fileRecord
	linkedJob *CopyJob
	walked    bool
	walkErr   error
	hasFiles  bool
	wipeDirOn bool
}

// ReadJobdir reads the contents of job.PathIn or job.PathOut
// and returns a checkDir struct.
// Walk error logged to checkDir.walkErr
func ReadJobdir(job *CopyJob, pathIn bool) *checkDir {
	var dir string
	if pathIn {
		dir = job.PathIn
	} else {
		dir = job.PathOut
	}
	cd := ReadDir(dir)
	return cd
}

// ReadDir reads the contents of dirpath and returns a checkDir struct.
// Walk error logged to checkDir.walkErr
func ReadDir(dirpath string) *checkDir {
	dirpath = MakeAbs(dirpath)
	cd := &checkDir{dirpath: dirpath, record: fileRecord{files: make(map[string]*filedata)}}
	err := fs.WalkDir(os.DirFS(dirpath), ".", func(p string, d fs.DirEntry, e error) error {
		if p == "." {
			return nil
		}
		_, ew := cd.record.newRecord(Joinpath(dirpath, p), p, d.IsDir())
		if ew != nil {
			return ew
		}
		return nil
	})
	cd.walked = true
	cd.walkErr = err
	return cd
}

func (C *checkDir) deleteDir() error {
	if C.wipeDirOn {
		return os.RemoveAll(C.dirpath)
	}
	return fmt.Errorf("wipeDirOn is false")
}

func (C *checkDir) String() string {
	out := uout.NewOutf("Detail dir: '%s'", C.dirpath)
	out.IndR()
	if C.linkedJob != nil {
		out.IfV(C.linkedJob.jobRan, "(job-linked: ran)", "(job-linked: not ran)")
	}
	out.IfF(C.walked, "Dir read - %d existing entries", "%s", len(C.record.files), "Dir not read")
	out.IndR().LV(C.record.files)
	out.IndL()
	if C.walkErr != nil {
		out.F("Walk error: %v", C.walkErr)
	}
	return out.String()
}
