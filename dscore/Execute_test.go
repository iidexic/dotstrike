package dscore

import (
	"fmt"
	"maps"
	"testing"

	pops "iidexic.dotstrike/pathops"
)

func (J *jobProcessor) testSetupAndDryRun(abortOnError bool) error {
	// if confirmHaveRuntimeConfig && (J.runtimeConfig == nil || len(J.runtimeConfig) == 0) {
	// 	return fmt.Errorf("Run terminated: No runtime config (confirmHaveRuntimeConfig on)")
	// }
	J.runtimeConfig[BoolNoFiles] = true
	var runErr error
	for i := range J.specs {
		J.specs[i].applyConfigsPrioritized(gd.data.Prefs.Bools, J.runtimeConfig)
		J.specs[i].group = pops.Copier().NewJobGroup(J.specs[i].groupExport())
		e := J.specs[i].group.RunAll(abortOnError)
		if e != nil {
			if abortOnError {
				return fmt.Errorf("Problem with run of group %s: %w", J.specs[i].Alias, e)
			}
			if runErr == nil {
				runErr = fmt.Errorf("Group Fails: (GRP-%s - %w)", J.specs[i].Alias, e)
			} else {
				runErr = fmt.Errorf("%w (GRP-%s - %w})", runErr, J.specs[i].Alias, e)
			}
		}

	}
	return runErr
}
func testConfig() map[ConfigOption]bool {
	m := make(map[ConfigOption]bool, 5)
	m[BoolIgnoreRepo] = true
	m[BoolKillGlobalTarget] = true
	m[BoolRootSubdir] = true
	m[BoolSeparateSources] = true
	m[BoolNoFiles] = true
	return m
}

func TestTempAssign(t *testing.T) {
	p := prefs{Bools: map[ConfigOption]bool{BoolIgnoreRepo: false, BoolIgnoreHidden: true, BoolUseGlobalTarget: false}}
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

func testSpec() *Spec {
	return &Spec{
		Alias: "@TEST_SPEC", Overrides: prefs{Bools: make(map[ConfigOption]bool)}, Ctype: specComponent,
		Sources: []pathComponent{{Path: "d:/coding/exampleFiles/imagesets", Ctype: sourceComponent}},
		Targets: []pathComponent{{Path: `d:\coding\exampleFiles\OUTPUT\ImageSets`, Ctype: targetComponent}},
	}
}

func fullTestSetupLazy(t *testing.T, useSelected bool) *jobProcessor {
	td := initForTest(t)
	jm := JobManager()
	jm.RuntimeConfigure(testConfig())
	if useSelected {
		jm.AddSpecs(td.SelectedSpec())
	} else {
		jm.AddSpecs(testSpec())
	}
	return jm
}

func TestJobSpecConfig(t *testing.T) {
	td := initForTest(t)
	jm := JobManager()
	tcfg := testConfig()
	t.Logf("user data: %+v", td)

	t.Logf("Selected: %s", td.SelectedSpec().Detail())
	jm.AddSpecs(td.SelectedSpec())

	t.Logf("Apply testConfig: %v", tcfg)
	jm.RuntimeConfigure(tcfg)
	t.Log(jm.WriteJobDetail())
	for k, v := range jm.runtimeConfig {
		val, ok := jm.runtimeConfig[k]
		if !ok {
			t.Errorf("runtime config missing key `%s` in testConfig.", k.String())
		} else if ok && val != v {
			t.Errorf("config failed to write for `%s`: expected %t, got %t", k.String(), v, val)
		} else {
			t.Logf("%s: %t == %t (test)", k.String(), v, val)
		}

	}
	maps.DeleteFunc(jm.runtimeConfig,
		func(k ConfigOption, v bool) bool { _, ok := tcfg[k]; return ok })
	t.Logf("non-test options: %v", jm.runtimeConfig)

	//-- run: --
	// e := jm.testSetupAndDryRun(true)
	// if e != nil {
	// 	t.Errorf("|testSetupDryRunError|\n%v", e)
	// } else {
	// 	t.Log(pops.Copier().GlobalOut)
	// }
}

func TestManagerToCopier(t *testing.T) {
	mgr := fullTestSetupLazy(t, false)
	t.Logf("%+v", mgr)
	// groups := pops.Copier().JobGroups
	// _ = groups
	t.Logf("Copier Detail:\n%s", pops.Copier().Detail())
}
