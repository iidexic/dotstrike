package dscore

import (
	"testing"

	pops "iidexic.dotstrike/pathops"
)

func TestAddComponent(t *testing.T) {
	S := Spec{Alias: "test", Ctype: specComponent, Sources: make([]PathComponent, 1), Targets: make([]PathComponent, 1)}
	// Test Paths
	td1 := "d:/coding/examplefiles/TEST-1"
	td2 := "d:/coding/examplefiles/TEST-2"
	S.AddSource(td1)
	S.CheckAddPath(td2, false)

	t.Logf("Spec Sources: %v", S.Sources)
	t.Logf("Spec Targets: %v", S.Targets)

	added := S.CheckAddPath(td1, false)
	t.Logf("Attempt to add %s as TARGET: added == %v", td1, added)
	if added {
		t.Errorf("Spec added existing path as target")
	} else {
		t.Errorf("Spec did not add existing path as TARGET")
	}
	e := S.AddSource(td2)
	t.Logf("Attempt to add %s as SOURCE", td2)

	if e != nil {
		t.Log(e)
	} else {
		t.Errorf("Spec added existing path as target")
	}

}
func testIsPathChild(t *testing.T, S *Spec, path string) bool {

	t.Logf("before and after clean:\n%s\n%s", path, pops.CleanPath(path))
	t.Logf("after tilde expand(just to check):%s", pops.TildeExpand(path))
	path = pops.CleanPath(path)
	for i, src := range S.Sources {
		t.Logf("Check: %d == %s? -> %s", i, src.Path, path)
		if src.Alias == path || src.Path == pops.MakeAbs(path) ||
			src.Path == path || src.Path == pops.CleanPath(path) ||
			src.Path == pops.TildeExpand(path) {
			return true
		}
	}
	for i, tgt := range S.Targets {
		t.Logf("Check: %d == %s? -> %s", i, tgt.Path, path)
		if tgt.Alias == path || tgt.Path == pops.MakeAbs(path) ||
			tgt.Path == path || tgt.Path == pops.CleanPath(path) ||
			tgt.Path == pops.TildeExpand(path) {
			return true
		}
	}
	return false
}

func TestIsPathChild(t *testing.T) {
	S := &Spec{Alias: "test", Ctype: specComponent, Sources: make([]PathComponent, 1), Targets: make([]PathComponent, 1)}
	// Test Paths
	td1 := "d:/coding/examplefiles/TEST-1"
	td2 := "d:/coding/examplefiles/TEST-2"
	// add
	S.AddSource(td1)
	S.CheckAddPath(td2, false)

	t.Logf("Spec Sources: %v", S.Sources)
	t.Logf("Spec Targets: %v", S.Targets)

	if testIsPathChild(t, S, td1) {
		t.Log("Path found in spec")
	}
	e := S.AddSource(td2)
	t.Logf("Attempt to add %s as SOURCE", td2)

	if e != nil {
		t.Log(e)
	} else {
		t.Errorf("Spec added existing path as target")
	}
}

func TestDeleteIfChildTilde(t *testing.T) {
	testConfig()
	S := &Spec{Alias: "test", Ctype: specComponent, Sources: make([]PathComponent, 1), Targets: make([]PathComponent, 1)}
	S.AddSource("~")
	t.Logf("1. Spec: %v", S.DetailFlat())
	S.CheckAddPath("~", true)
	t.Logf("2. Spec: %v", S.DetailFlat())
	S.DeleteIfChild("~", false, true)
	t.Logf("3. Spec %v", S.DetailFlat())
}
