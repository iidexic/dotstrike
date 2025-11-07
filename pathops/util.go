package pops

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type dirlog struct {
	path                         string
	job                          *CopyJob
	readPre, readPostDone, Clear bool
	pre, post                    []filedetail
	ePre, ePost                  error
}

type filedetail struct {
	path string
	size int64
	mod  time.Time
	dir  bool
}

func newDirLog(job *CopyJob) (*dirlog, error) {
	absDir := job.PathOut

	if job.BPrefs.IsOn(bRootSubdir) {

	}
	dl := dirlog{path: absDir, pre: make([]filedetail, 0), post: make([]filedetail, 0)}
	dl.job = job
	e := fs.WalkDir(os.DirFS(absDir), ".", func(p string, d fs.DirEntry, e error) error {
		i := len(dl.pre)
		dl.pre = append(dl.pre, filedetail{path: Joinpath(absDir, p), dir: d.IsDir()})
		if !dl.pre[i].dir {
			info, e := d.Info()
			if e != nil {
				return nil
			}
			dl.pre[i].size = info.Size()
			dl.pre[i].mod = info.ModTime()
		}

		return nil
	})
	dl.readPre = true
	if e != nil {
		if dl.ePre != nil {
			dl.ePre = fmt.Errorf("%w, %w", dl.ePre, e)
		} else {
			dl.ePre = fmt.Errorf("WalkDir Error: %w", e)
		}
		return &dl, e // I think this is unnecessary
	}
	return &dl, nil
}

func (D *dirlog) _walk_dir_() {
	var fl []filedetail
	e := fs.WalkDir(os.DirFS(D.path), ".", func(p string, d fs.DirEntry, e error) error {
		i := len(fl)
		fl = append(fl, filedetail{path: Joinpath(D.path, p)})
		isdir := d.IsDir()
		fl[i].dir = isdir
		if !isdir {
			info, e := d.Info()
			if e != nil {
				return nil
			}
			fl[i].size = info.Size()
			fl[i].mod = info.ModTime()
		}
		return nil
	})
	if D.readPre && !D.readPostDone {
		D.post = fl
		D.readPostDone = true
		if e != nil {
			if D.ePost != nil {
				D.ePost = fmt.Errorf("%w, %w", D.ePost, e)
			} else {
				D.ePost = fmt.Errorf("WalkDir Error: %w", e)
			}
		}
	} else if !D.readPre {
		D.pre = fl
		D.readPre = true
		if e != nil {
			if D.ePre != nil {
				D.ePre = fmt.Errorf("%w, %w", D.ePre, e)
			} else {
				D.ePre = fmt.Errorf("WalkDir Error: %w", e)
			}
		}
	}
}

// TODO: (Hi) Don't need to do a read post; should have collected all details in main WalkDir
// TODO: ACTUALLY This should just be built into the CopyJob, no need to run two walkdirs.
func (D *dirlog) readPost() error {
	if !D.readPre {
		return fmt.Errorf("Pre-state not recorded")
	}
	D._walk_dir_()
	return nil
}

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
	if p1 == p2 {
		return true
	}
	return false
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

// Just checks equality between filename/ext and size
func looksLikeRawCopy(checkPath string, basefile fs.FileInfo) bool {
	cstat, err := os.Stat(checkPath)
	if err != nil {
		return false
	}
	return cstat.Name() == basefile.Name() && cstat.Size() == basefile.Size() && cstat.IsDir() == basefile.IsDir()
}
