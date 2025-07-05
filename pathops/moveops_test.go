package pops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func checkbasicpath(fpath string, dest string, t *testing.T) {
	if !IsBasicPath(fpath) {
		t.Logf("not basic path: '%s'", fpath)
		t.Fail()
	}
	if !IsBasicPath(dest) {
		t.Logf("not basic path: '%s'", dest)
		t.Fail()
	}
	fpdetail := DetailStatPath(fpath)
	ddetail := DetailStatPath(dest)

	t.Log(fpdetail)
	t.Log(ddetail)

}

func findProjdir(upath, projname string, t *testing.T) string {
	if !filepath.IsAbs(upath) {
		if filepath.Ext(upath) == "" {
			//prob dir
		}

		e := os.Chdir("..")
		if e != nil {
			t.Logf("Chdir error: %e", e)
			e = nil
		}
		uabs, e := filepath.Abs(upath)
		if e != nil {
			t.Logf("Get abspath from '%s' failed: %e", upath, e)
		} else {
			upath = uabs
		}
	}
	return upath
}

func testLilCopy(fpath string, dest string, t *testing.T) {

	//checkbasicpath(fpath, dest, t)
	filefrom, efrom := OpenExistingFile(fpath)
	if efrom != nil {
		t.Logf("error opening filefrom: %e", efrom)
	}
	if filefrom == nil {
		t.Error("filefrom: no file data")
	}
	defer filefrom.Close()
	dpath, e := filepath.Abs(dest)
	if e != nil {
		t.Error(e)
	}

	fileto, eto := MakeOpenFileF(dpath)
	t.Logf("file-to Fd: %v", fileto.Fd())
	t.Log(eto)
	if eto != nil {
		t.Log("Error MakeOpenFile: ")
		t.Error(eto)
	}
	defer fileto.Close()
	written, ecpy := io.Copy(fileto, filefrom)
	if ecpy != nil {
		t.Log("(first) Error Copying")
		t.Error(ecpy)
	}
	t.Log(fmt.Sprintf("bytes copied: %d", written))

}

func TestLilCopyGetToml(t *testing.T) {
	configtoml := HomeDirtyJoin(".config/dotstrike/dotstrikeData.toml")
	loctoml := findProjdir("./_xtra/dotstrikeData.toml", "dotstrike", t)
	testLilCopy(configtoml, loctoml, t)
}

func TestLilCopyPushToml(t *testing.T) {
	testLilCopy("./_xtra/dotstrikeData.toml", "~/.config/dotstrike/dotstrikeData.toml", t)
}

func TestStatHomePath(t *testing.T) {
	//paths := []string{"/~/", "\\~\\", `~\\`, `~~//`}
	hdir, errh := HomeJoin(".config/")
	if errh != nil {
		t.Errorf("HomeDir Error %v", errh)
	}
	info, e := os.Stat(hdir)
	infostr := fmt.Sprintf(`path:%s stat:%+v`, hdir, info)
	if e != nil {
		t.Logf("%s | error %v", infostr, e)
	} else {
		t.Log(infostr)
		t.Logf("dir? -> %t", info.IsDir())
	}
}

func TestLStatPaths(t *testing.T) {
	paths := []string{".", "./", `./`, ".\\", `.\\`, "~", `~`,
		"~/", `~/`, "~/.config", `~/.config`, "~//.config", `~//.config`,
		"~/.config/dotstrike/dotstrikeData.toml", "~//.config//dotstrike//dotstrikeData.toml",
		`~/.config/dotstrike/dotstrikeData.toml`, `c:\~\`, `d:\~\`}

	for i, p := range paths {
		info, e := os.Lstat(p)
		infostr := fmt.Sprintf(`%d) path:%s stat:%+v`, i, p, info)
		if e != nil {
			t.Logf("%s | error %v", infostr, e)
		} else {
			t.Log(infostr)
		}
	}
}
func TestStatPaths(t *testing.T) {
	paths := []string{".", "./", `./`, ".\\", `.\\`, "~", `~`,
		"~/", `~/`, "~/.config", `~/.config`, "~//.config", `~//.config`,
		"~/.config/dotstrike/dotstrikeData.toml", "~//.config//dotstrike//dotstrikeData.toml",
		`~/.config/dotstrike/dotstrikeData.toml`, `c:\~\`, `d:\~\`}

	for i, p := range paths {
		info, e := os.Stat(p)
		infostr := fmt.Sprintf(`%d) path:%s stat:%+v`, i, p, info)
		if e != nil {
			t.Logf("%s | error %v", infostr, e)
		} else {
			t.Log(infostr)
		}
	}
}
func TestPrintDebug(t *testing.T) {
	t.Log()
}
