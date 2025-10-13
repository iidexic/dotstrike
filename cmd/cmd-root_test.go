package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/uout"
)

func commandTestList() map[string]string {
	return map[string]string{
		"list":                         "User Specs:",
		"spec":                         "Selected spec:",
		"spec commandtest -y":          "",
		"spec commandtest --delete -y": "deleted",
	}

}
func CommandRef() map[string]*cobra.Command {
	cl := rootCmd.Commands()
	lookup := make(map[string]*cobra.Command, len(cl))
	for _, cmd := range cl {
		lookup[cmd.Name()] = cmd
	}
	return lookup
}

type tRunner struct {
	inputs  []string
	outputs []string
	errors  []error
}

func testRunSequence(inputs []string, t *testing.T) (*tRunner, error) {
	r := testCmdRunner(inputs)
	r.Execute()
	t.Logf("%v", r)
	var e error
	for _, err := range r.errors {
		if err != nil {
			if e == nil {
				e = fmt.Errorf("%w", err)
			} else {
				e = fmt.Errorf("%w, %w", e, err)
			}
		}

	}
	return r, e
}

func (R tRunner) String() string {
	out := uout.NewOut("Run Results")
	for i, o := range R.outputs {
		out.F("IN: '%s'", R.inputs[i])
		if e := R.errors[i]; e != nil {
			out.F("Error: %v", e)
		}
		out.F("OUT: %s", o)
		out.Sep()
	}
	return out.String()
}

func testCmdRunner(inputs []string) *tRunner {
	return &tRunner{inputs: inputs, outputs: make([]string, len(inputs)), errors: make([]error, len(inputs))}
}


func (R *tRunner) Execute() {
	for i, runarg := range R.inputs {
		R.outputs[i], R.errors[i] = testExec(rootCmd, runarg)
	}
}

func (R *tRunner) ExecuteLog(t *testing.T) {
	for i, runarg := range R.inputs {

		R.outputs[i], R.errors[i] = testExec(rootCmd, runarg)
		if R.errors[i] != nil {
			t.Logf("Executed cmd (%d): ERR: %s", i, R.errors[i].Error())
		} else {
			t.Logf("Executed cmd (%d)", i)
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
	bin := bytes.NewReader([]byte(args))
	bout := bytes.NewBufferString("")
	cmd.SetIn(bin)
	cmd.SetArgs(strings.Split(args, " "))
	cmd.SetOut(bout)
	e := cmd.Execute()
	return bout.String(), e
}

func testRoot(args string) (string, error) { return testExec(runCmd, args) }

func testRootSl(args string) ([]string, error) { return testCmdLines(runCmd, args) }

func TestTestCommand(t *testing.T) {
	execArgs := "cfg --global"
	sout, e2 := testExec(rootCmd, execArgs)
	if e2 != nil {
		t.Errorf("Single-String Execute Error: %s", e2.Error())
	}
	ez := uout.NewOut("Output(lines):")
	ez.WipeOnOutput(true)
	ez.V("Output(string):")
	ez.V(sout)
	out, e := testCmdLines(rootCmd, execArgs)
	if e != nil {
		t.Errorf("Lines Execute Error: %s", e.Error())
	}
	ez.ILV(out)
	t.Logf("Results Run '%s':\n%s", execArgs, ez.String())
	ea2 := "spec"
	out, e = testCmdLines(rootCmd, ea2)
	if e != nil {
		t.Errorf("Lines Execute Error: %s", e.Error())
	}
	ez.V("Output(lines):")
	ez.ILV(out)
	sout, e2 = testExec(rootCmd, ea2)
	if e2 != nil {
		t.Errorf("Single-String Execute Error: %s", e2.Error())
	}
	ez.V("Output(string):")
	ez.V(sout)
	t.Logf("Results run '%s' :\n%s", ea2, ez.String())
}

func TestTestCommandRunner(t *testing.T) {
	execArgs := []string{"cfg --global"}
	run := testCmdRunner(execArgs)
	t.Logf("run inputs: %v", run.inputs)
	run.Execute()
	t.Log("Output(lines):")
	for i, s := range run.outputs {
		t.Logf("%d) %s", i, s)
	}
	for i, e := range run.errors {
		if e != nil {
			t.Errorf("Execute Error#%d: %v", i, e)
		}
	}
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
func TestStructure(t *testing.T) {
	ins := []string{"spec srcTest", "src C:/secret/bringo", "src"}
	_, e := testRunSequence(ins, t)
	if e != nil {
		t.Errorf("Failures during run sequence")
	}

}

// TODO:(mid) finish this guyy
func TestFeatureset(t *testing.T) {
	/* When will confirmation be required:
	- [1] shouldn't but does - fix
	- [3]
	- [4]
	*/
	in := []string{
		//[0] make 2 spec (good)
		"spec test-sound test-svg",
		// Prep img folder test
		// [1] make img spec with 2 src, 1 tgt
		"spec test-img --src='d:/coding/exampleFiles/imagesets/svg-x-circle,d:/coding/exampleFiles/imagesets/svg_circle' --tgt=d:/coding/exampleFiles/OUTPUT/images -y",
		// [2] add last src
		"src d:/coding/exampleFiles/imagesets/svg_png",
		// [3] test delete 2 specs
		"spec test-sound test-svg --delete -y",
		// [4] init both audio test specs
		"spec test-audio test-audiodirs --src=d:/coding/exampleFiles/audio -y",
		// [5] add tgt to test-audio
		"tgt d:/coding/exampleFiles/OUTPUT/audio --ignore=*.mp3",
		// [6] select test-audiodirs
		"sel iodir",
		// [7] set cfg for test-audiodirs
		"cfg dry makealldirs",
		// [8] add tgt to test-audiodirs
		"tgt d:/coding/exampleFiles/OUTPUT/audio-structure",
		// [9] list (check output correct)
		"list",
		//cleanup
		"spec test-img test-audio test-audiodirs --delete -y",
	}
	run := testCmdRunner(in)
	run.ExecuteLog(t)
	t.Logf("%v", *run)

	// for i, o := range run.outputs {
	// 	if e := run.errors[i]; e != nil {
	// 		t.Errorf("Run %d Error: %v", i, e)
	// 	}
	// 	t.Logf(`Run %d
	// INPUT: %s
	// OUTPUT: %s`, i, in[i], o)
	// }

}
