package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	pops "iidexic.dotstrike/pathops"
)

func printFlagInfo(t *testing.T, cmd *cobra.Command) {
	t.Logf("LOOKING AT: %s", cmd.Name())
	t.Log("0.PRINT COMMAND-------")
	t.Logf("%+v", *cmd)
	t.Log("1.HELP----------------")

	t.Log(cmd.Help())
	t.Log("2.USELINE-------------")
	t.Log(cmd.UseLine())
	t.Log("3.USAGE---------------")
	t.Log(cmd.Usage())
	t.Log("4.FLAGS---------------")
}

func testPrintFlagsVisit(t *testing.T, cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(p *pflag.Flag) {
		t.Logf("\n============[ Flag: %s ]============|", p.Name)
		t.Logf("Usage: %s", p.Usage)
		t.Logf("flag.Changed = %t", p.Changed)
		t.Logf("full flag:\n%+v", p)
	})
}

func TestRunFlag(t *testing.T) {
	printFlagInfo(t, runCmd)
	testPrintFlagsVisit(t, runCmd)
}

func TestRunDetail(t *testing.T) {

	out, e := testCmdLines(rootCmd, "run --setup-only-debug")
	if e != nil {
		t.Errorf("testRoot execute error: %v", e)
	}

	t.Logf("len out = %d", len(out))
	t.Logf("out:\n%s", strings.Join(out, "\n"))

	t.Logf("MAIN RUN DETAIL:")
	t.Log(detailMainRun())
	t.Log(" COPIER DETAIL:")
	t.Log(pops.Copier().Detail())

}
