package cmd

import "testing"

func TestSrcCmdBasicAll(t *testing.T) {
	ins := []string{"spec srcTest", "src C:/secret/bringo", "src"}
	_, e := testRunSequence(ins, t)
	if e != nil {
		t.Errorf("Failures during run sequence")
	}

}

func TestTildeError(t *testing.T) {
	ins := []string{"spec srcTest --status-report", "src ~ --status-report", "src ~ --delete --status-report"}
	runner := testCmdRunner(ins)
	for !runner.Done() {
		runner.ExecuteNextLog(t)
		if runner.errors[runner.runIndex-1] != nil {
			t.Errorf("ERROR")
		}
	}
	t.Logf("%v", *runner)
}
