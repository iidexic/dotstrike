package cmd

import "testing"

func TestLearnStructCreationBehavior(t *testing.T) {
	run := tRunner{}
	t.Log(len(run.cmdList))
	t.Log(len(run.inputs))
	t.Logf("Nil? -> %t", run.inputs == nil)
	t.Log("append?")
	run.inputs = append(run.inputs, "poo", "pee", "fart", "shit", "cum")
	t.Logf("ok did it work. len now: %d", len(run.inputs))
	t.Log(len(run.outputs))
	t.Log(len(run.errors))
	srand := make([]string, 0)
	t.Logf("randmake: len is %d, is nil -> %t", len(srand), srand == nil)
	srandp := make([]string, 0, 1000)
	t.Logf("randmake w/1000 prep: len is %d, is nil -> %t", len(srandp), srandp == nil)
}

func TestSrcCmdBasicAll(t *testing.T) {
	run := tRunner{}
	tot := run.addInputs(
		"spec srcTest",
		"src C:/secret/bringo",
		"src",
	)
	for i, s := range run.inputs {
		t.Logf("in #%d: '%s'", i, s)
	}
	if tot != 3 { /*magicnumber*/
		t.Errorf("addInputs failed or magic# didn't get changed but the test inputs did :)")
	}
	run.Execute(false)
	for i, e := range run.errors {
		if e != nil {
			t.Errorf("Error from input [%d]:\nIN:%s\nOUT:%s\nERROR:%s", i,
				run.inputs[i], run.outputs[i], run.errors[i].Error())
		} else {
			t.Logf("Input [%d] good\nIN:%s\nOUT:%s", i, run.inputs[i], run.outputs[i])
		}
	}
}
