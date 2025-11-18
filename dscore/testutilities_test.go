package dscore

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	pops "iidexic.dotstrike/pathops"
	"iidexic.dotstrike/uout"
)

// ── Decoding ────────────────────────────────────────────────────────

func initForTest(t *testing.T) *globalModify {
	LoadGlobals()  // decode globals
	InitTempData() // load tempdata struct with globals details
	temp := TempData()
	if temp == nil {
		t.Error("Temp Data not initialized")
	}
	return temp
}

// ── Encoding ────────────────────────────────────────────────────────

func encodeTomltesting(path string, data *globalData) error {
	file, e := pops.OpenFileRW(pops.CleanPath(path))
	if file != nil {
		defer file.Close()
	}
	if e != nil || file == nil {
		return fmt.Errorf("Error opening toml for write: %w", e)
	}
	encode := toml.NewEncoder(file)
	e = encode.Encode(*data)
	if e != nil {
		return e
	} else {
		return nil
	}
}
func encodeToBuffer(data *globalData) (bytes.Buffer, error) {
	buf := bytes.Buffer{}
	e := toml.NewEncoder(&buf).Encode(*data)
	if e != nil {
		return buf, e
	}
	return buf, nil
}

// ── Error logging ───────────────────────────────────────────────────
func tLogErr(context string, e error, t *testing.T) {
	if e != nil {
		t.Logf("%s: %s", context, e.Error())
	}
}

func dumpGlobalLog(t *testing.T) {
	out := uout.NewOut("[ Global Log ]")
	out.ILV(gd.GlobalMessage)
	t.Log(out.String())
}

// same as globalModify_test.go encodeTomltesting
// NOTE: Needed?
func encodeTestfile(path string, data *globalData) error {
	file, e := pops.OpenFileRW(pops.CleanPath(path))
	if file != nil {
		defer file.Close()
	}
	if e != nil || file == nil {
		return fmt.Errorf("Error opening toml for write: %w", e)
	}
	encode := toml.NewEncoder(file)
	e = encode.Encode(*data)
	if e != nil {
		return e
	} else {
		return nil
	}
}
