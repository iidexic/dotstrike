package dscore

import (
	"testing"

	"iidexic.dotstrike/uout"
)

func initForTest(t *testing.T) *globalModify {
	LoadGlobals()
	InitTempData()
	temp := TempData()
	if temp == nil {
		t.Error("Temp Data not initialized")
	}
	return temp
}

func dumpGlobalLog(t *testing.T) {
	out := uout.NewOut("[ Global Log ]")
	out.ILV(gd.GlobalMessage)
	t.Log(out.String())
}
