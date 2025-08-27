package cmd

import (
	"bytes"
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

func TestRunAllCommands(t *testing.T) {
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
			t.Logf("IN:%s OUT:%s", k, v)
		}

	}

}
