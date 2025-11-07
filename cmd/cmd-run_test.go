package cmd

import (
	"strings"
	"testing"

	pops "iidexic.dotstrike/pathops"
)

func TestRunDetail(t *testing.T) {
	out, e := testCmdLines(rootCmd, "run --setup-only-debug")
	if e != nil {
		t.Errorf("execute error: %v", e)
	}

	t.Logf("len out = %d", len(out))
	t.Logf("out:\n%s", strings.Join(out, "\n"))

	t.Logf("MAIN RUN DETAIL:")
	t.Log(detailMainRun())
	t.Log(" COPIER DETAIL:")
	t.Log(pops.Copier().Detail())

}

func TestRunMultiSource(t *testing.T) {
	testseq := []string{
		"spec test-run-multi --src=d:/coding/exampleFiles/OUTPUT/images,d:/coding/exampleFiles/OUTPUT/audio --tgt=d:/coding/exampleFiles/OUTPUT/multi",
		"run"}
	runner := testCmdRunner(testseq)
	for !runner.Done() {
		runner.ExecuteNextLog(t)
		t.Logf("run %d in - %s", runner.runIndex, runner.inputs[runner.runIndex-1])
		t.Logf("run %d out - %s", runner.runIndex, runner.outputs[runner.runIndex-1])

	}
}
