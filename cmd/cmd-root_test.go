package cmd

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func commandTestList() map[string]string {
	return map[string]string{
		"list":                         "User Specs:",
		"spec":                         "Selected spec:",
		"spec commandtest -y":          "",
		"spec commandtest --delete -y": "deleted",
	}

}

type tRunner struct {
	cmdList []*cobra.Command
	inputs  []string
	outputs []string
	errors  []error
}

func (R *tRunner) addCommands(cmd ...*cobra.Command) int {
	R.cmdList = append(R.cmdList, cmd...)
	return len(R.cmdList)
}

func (R *tRunner) addInputs(in ...string) int {
	// if R.inputs == nil { // impossible
	// 	R.inputs = make([]string, len(in))
	// 	copy(R.inputs, in)
	// 	return len(R.inputs)
	// }
	R.inputs = append(R.inputs, in...)
	return len(R.inputs)
}

// TODO: uh make this work
func (R *tRunner) Execute(useCommands bool) {
	lin, lcmd := len(R.inputs), len(R.cmdList)
	R.outputs = make([]string, max(lin, lcmd))
	R.errors = make([]error, max(lin, lcmd))
	if useCommands {
		if ll, lc := len(R.inputs), len(R.cmdList); ll < lc {
			_ = slices.Grow(R.inputs, lc-ll)
		}
		for i := range R.cmdList {
			R.outputs[i], R.errors[i] = testExec(R.cmdList[i], R.inputs[i])
		}
	} else {
		for i, runarg := range R.inputs {
			R.outputs[i], R.errors[i] = testExec(runCmd, runarg)
		}
	}
}

func containsSubstring(text, sub string) bool {
	text = strings.ToLower(text)
	sub = strings.ToLower(sub)
	return strings.Contains(text, sub)
}

func testCmdLines(cmd *cobra.Command, args string) ([]string, error) {
	bout := bytes.NewBufferString("")
	cmd.SetArgs(strings.Split(args, " "))
	cmd.SetOut(bout)
	e := cmd.Execute()
	return strings.Split(bout.String(), "\n"), e
}

func testExec(cmd *cobra.Command, args string) (string, error) {
	//bin := bytes.NewReader([]byte(input))
	bout := bytes.NewBufferString("")
	cmd.SetArgs(strings.Split(args, " "))
	cmd.SetOut(bout)
	e := cmd.Execute()
	return bout.String(), e
}

func testRoot(args string) (string, error) { return testExec(runCmd, args) }

func testRootSl(args string) ([]string, error) { return testCmdLines(runCmd, args) }

func TestTestCommand(t *testing.T) {
	execArgs := "cfg --global"
	out, e := testCmdLines(rootCmd, execArgs)
	if e != nil {
		t.Errorf("Lines Execute Error: %s", e.Error())
	}
	t.Log("Output(lines):")
	for i, s := range out {
		t.Logf("%d) %s", i, s)
	}
	sout, e2 := testExec(rootCmd, execArgs)
	if e2 != nil {
		t.Errorf("Single-String Execute Error: %s", e2.Error())
	}
	t.Log("Output(string):")
	t.Log(sout)
}

func TestRunTestSequence(t *testing.T) {
	ct := commandTestList()
	for k, v := range ct {
		out, e := testExec(rootCmd, k)
		if e != nil {
			t.Errorf("fail from running `%s`\nfailure:%s", k, e.Error())
		}
		if !containsSubstring(out, v) {
			t.Errorf(`verifying string not found in output
args passed: "%s", verifying string: "%s"
Output:
"%s"`, k, v, out)
		} else {
			t.Logf("in:%s\n----------------------\nout:%s\n----------------------\n(validation out:%s)", k, out, v)
		}

	}

}
func runSequential(runargs ...string) ([]string, error) {
	var err error
	output := make([]string, len(runargs))
	for i, a := range runargs {
		s, e := testRoot(a)
		if e != nil {
			if err == nil {
				err = fmt.Errorf("cmd errors: [%d] %w,", i, e)
			} else {
				err = fmt.Errorf("%w [%d] %w, ", err, i, e)
			}
		}
		output[i] = s

	}
	return output, err
}

// TODO:(mid) finish this guyy
func TestFeatureset(t *testing.T) {
	out, err := runSequential(
		"spec test-audio test-svg",                          // make multiple spec
		"spec test-imagesets",                               // make  spec
		"src d:/coding/exampleFiles/imagesets/svg-sizediff", // add src
		"src d:/coding/exampleFiles/imagesets/svg-x-circle",
		"src d:/coding/exampleFiles/imagesets/svg_circle d:/coding/exampleFiles/imagesets/svg_png", //add 2 src
		"tgt d:/coding/exampleFiles/OUTPUT/images",
		"cfg ",                     //WARN:notdone
		"spec test-audio --delete", //test deletes
		"spec  test-svg --delete",
		// single line new spec with inline paths
		`spec t-audio --src='d:/coding/exampleFiles/audio' --tgt=='d:/coding/exampleFiles/OUTPUT/audio'`,
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

}
