package dscore

import "testing"

func TestTempAssign(t *testing.T) {
	p := prefs{KeepRepo: true, KeepHidden: true, GlobalTarget: false}
	//1. run init
	temp := initForTest(t)
	//2. get selected spec
	spec := temp.SelectedSpec()
	t.Log(spec.Detail())
	if !temp.Modified {
		spec.Overrides = p
		spec.OverrideOn = true
		t.Log("After Modifying:")
		t.Log(spec.Detail())
	} else {
		t.Error("tempdata is marked as modified for some reason")
	}

}
