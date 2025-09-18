package pops

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"iidexic.dotstrike/config"
	"iidexic.dotstrike/uout"
)

type dirRecord map[string]pathDetail

type pathDetail struct {
	//job                            *CopyJob
	exists, made, isDir            bool
	sizeOrigin, sizePre, sizeFinal int64
	modPre, modFinal               time.Time
}

// CopyJob prepares and executes the copy of all contents of PathIn to PathOut
type CopyJob struct {
	PathIn, PathOut string     // Root of copy source and destination (*or destination parent)
	parentPathOut   string     // unused. populated on run if JobSettings.makeRootSubdir = true
	fstack          []filecopy // record of files copied
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
		d += fmt.Sprintf("| ran (%d file, %d newdir, %d errors)", len(J.fstack), J.DirsMade(), len(J.OpErrors))
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
			len(J.fstack), J.CopyPercent()*100, len(J.newDirs), J.DirsMade())
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
	out := uout.NewOut("Files")
	out.IndR()
	out.ILV(J.fstack)
	return out.String()
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
	for _, f := range J.fstack {
		read += f.inSize
		copied += f.outSize
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

func (J *CopyJob) logPathError(perr PathError) {
	J.OpErrors = append(J.OpErrors, perr)
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

func (J *CopyJob) addFile(relpath string, inSize, outSize int64) {
	J.fstack = append(J.fstack, filecopy{relpath: relpath, inSize: inSize, outSize: outSize})
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

// stripRoot removes CopyJob.PathIn from path provided for construction of destination path
// structure/intent of CopyJob requires J.PathIn to be a prefix in rpath.
// As such, if an error is encountered, stripRoot panics
func (J *CopyJob) stripRoot(p string) string {

	relp, e := filepath.Rel(J.PathIn, p)
	if e != nil {
		panic(fmt.Errorf("stripRoot(%s) error: %v", p, e))
	}
	return relp

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
