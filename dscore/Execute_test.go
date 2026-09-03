package dscore

import (
	"maps"
	"testing"

	"iidexic.dotstrike/config"
)

func testConfig() map[ConfigOption]bool {
	m := make(map[ConfigOption]bool, 5)
	m[BoolIgnoreRepo] = true
	m[BoolKillGlobalTarget] = true
	m[BoolRootSubdir] = true
	//m[BoolSeparateSources] = true
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
		Sources: []PathComponent{{Path: "d:/coding/exampleFiles/imagesets", Ctype: sourceComponent}},
		Targets: []PathComponent{{Path: `d:\coding\exampleFiles\OUTPUT\ImageSets`, Ctype: targetComponent}},
	}
}

// Runs initForTest then adds 1 or 2 specs; adds testSpec() and if useSelected then adds selected spec
func fullTestSetupLazy(t *testing.T, useSelected bool) *jobProcessor {
	td := initForTest(t)
	jm := JobManager()
	jm.RuntimeConfigure(testConfig())
	if useSelected {
		jm.AddSpecs(td.SelectedSpec())
	}
	jm.AddSpecs(testSpec())
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
	mgr.SetupOnly()
	t.Logf("%+v", mgr)
	// groups := pops.Copier().JobGroups
	// _ = groups
	//copier := pops.Copier()
	//cpstr := copier.String()
	//t.Logf("Copier Detail:\n%s", cpstr)
	for _, spec := range mgr.specs {
		t.Logf("SPEC CONFIG:\n%+v", spec.config)
		if spec.group != nil {
			t.Logf("Spec %s -> Group %s", spec.Alias, spec.group.Name())
			groupConfig := spec.group.Config()
			t.Logf("GROUP CONFIG:\n%+v", groupConfig)
			if config.ConfigsMatch(spec.config, groupConfig) {
				t.Logf("Group Config Matches Spec Config")
			} else {
				t.Errorf("Group Config does not match Spec Config")
			}
			for _, job := range spec.group.CopyJobs() {
				t.Logf("JOB CONFIG:\n%+v", job.BPrefs)
				if !config.ConfigsMatch(groupConfig, job.BPrefs) {
					t.Errorf("Job Config does not match Job Group Config")
				}
			}
		}
	}
}

func TestRunCopy(t *testing.T) {
	mgr := fullTestSetupLazy(t, false)
	_ = mgr
}
