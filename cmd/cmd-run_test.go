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
