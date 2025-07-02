package pops

import (
	"fmt"
	"io"
	"os"
	"path"
	"testing"
)

func testhome(suffix string, t *testing.T) string {
	h, e := HomeJoin(suffix)
	if e != nil {
		t.Log("HomeDir Error:", e)
	}
	return path.Join(h, suffix)
}

func testLilCopy(fpath string, dest string, t *testing.T) {
	if !IsBasicPath(fpath) {
		t.Logf("not basic path: '%s'", fpath)
		t.Fail()
	}
	if !IsBasicPath(dest) {
		t.Logf("not basic path: '%s'", dest)
		t.Fail()
	}
	fpdat, efrom := os.Stat(fpath)
	ce(efrom)
	ddat, efrom := os.Stat(dest)

	t.Logf("fp '%s' stat: %+v", fpath, fpdat)
	t.Logf("dest '%s' stat: %+v", dest, ddat)

	filefrom, efrom := OpenExistingFile(fpath)
	if efrom != nil {
		t.Log(efrom)
	}
	if filefrom == nil {
		t.Error("filefrom: no file data")
	}
	defer filefrom.Close()
	fileto, eto := MakeOpenFileF(dest)
	if eto != nil {
		t.Error(eto)
	}
	defer fileto.Close()
	written, ecpy := io.Copy(filefrom, fileto)
	if ecpy != nil {
		t.Log("(first) Error Copying")
		t.Error(ecpy)
	}
	t.Log(fmt.Sprintf("bytes copied: %d", written))

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

func TestLilCopyGetToml(t *testing.T) {
	testLilCopy("~/.config/dotstrike/dotstrikeData.toml", "./_xtra/dotstrikeData.toml", t)
}

func TestLilCopyPushToml(t *testing.T) {
	testLilCopy("./_xtra/dotstrikeData.toml", "~/.config/dotstrike/dotstrikeData.toml", t)
}
