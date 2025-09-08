package cmd

import (
	"bytes"
	"fmt"
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
func containsSubstring(text, sub string) bool {
	text = strings.ToLower(text)
	sub = strings.ToLower(sub)
	return strings.Contains(text, sub)
}

func testCommand(cmd *cobra.Command, args string) ([]string, error) {
	bout := bytes.NewBufferString("")
	cmd.SetArgs(strings.Split(args, " "))
	cmd.SetOut(bout)
	e := cmd.Execute()
	return strings.Split(bout.String(), "\n"), e
}

func testCmdString(cmd *cobra.Command, args string) (string, error) {
	//bin := bytes.NewReader([]byte(input))
	bout := bytes.NewBufferString("")
	cmd.SetArgs(strings.Split(args, " "))
	cmd.SetOut(bout)
	e := cmd.Execute()
	return bout.String(), e
}

func testRoot(args string) (string, error) { return testCmdString(runCmd, args) }

func testRootSl(args string) ([]string, error) { return testCommand(runCmd, args) }

func TestTestCommand(t *testing.T) {
	out, e := testCommand(rootCmd, "list")
	if e != nil {
		t.Errorf("Execute Error: %s", e.Error())
	}
	t.Log("Output:")
	for i, s := range out {
		t.Logf("[%d] %s", i, s)
	}
}

func TestRunTestSequence(t *testing.T) {
	ct := commandTestList()
	for k, v := range ct {
		out, e := testCmdString(rootCmd, k)
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
		"spec test-audio",
		"spec test-svg",
		"spec test-imagesets",
		"src d:/coding/exampleFiles/imagesets/svg-sizediff ",
		"src d:/coding/exampleFiles/imagesets/svg-x-circle",
		"src d:/coding/exampleFiles/imagesets/svg_circle",
		"src d:/coding/exampleFiles/imagesets/svg_png",
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

}
