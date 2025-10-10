package cmd

import "testing"

func TestLearnStructCreationBehavior(t *testing.T) {
	run := tRunner{}
	t.Log(len(run.inputs))
	t.Logf("Nil? -> %t", run.inputs == nil)
	t.Log("append?")
	run.inputs = append(run.inputs, "get", "list", "cfg", "check")
	t.Logf("ok did it work. len now: %d", len(run.inputs))
	run.Execute()
	t.Log(len(run.outputs))
	t.Log(len(run.errors))
	srand := make([]string, 0)
	t.Logf("randmake: len is %d, is nil -> %t", len(srand), srand == nil)
	srandp := make([]string, 0, 1000)
	t.Logf("randmake w/1000 prep: len is %d, is nil -> %t", len(srandp), srandp == nil)
}

func TestSrcCmdBasicAll(t *testing.T) {
	ins := []string{"spec srcTest", "src C:/secret/bringo", "src"}
	_, e := testRunSequence(ins, t)
	if e != nil {
		t.Errorf("Failures during run sequence")
	}

}
