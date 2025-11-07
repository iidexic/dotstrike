package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	"iidexic.dotstrike/uout"
)

type tRunner struct {
	inputs   []string
	outputs  []string
	errors   []error
	runIndex int
	verbose  bool
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
		R.runIndex++
		testSetFlagsDefault()
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
		// oh boy
		if R.verbose {
			t.Logf("TEMPDATA BEFORE RESET:\n%s", dscore.TempData().DetailFlat())
		}
		testSetFlagsDefault()
		cycleCoreForTest()
		//
		if R.verbose {
			t.Logf("TEMPDATA:\n%s", dscore.TempData().DetailFlat())
		}
		R.runIndex++
	}
}
func (R *tRunner) ExecuteNext() {
	if R.runIndex < len(R.inputs) {
		R.outputs[R.runIndex], R.errors[R.runIndex] = testExec(rootCmd, R.inputs[R.runIndex])
		R.runIndex++
	}
}

func (R *tRunner) ExecuteNextLog(t *testing.T) {
	if R.runIndex < len(R.inputs) {
		R.outputs[R.runIndex], R.errors[R.runIndex] = testExec(rootCmd, R.inputs[R.runIndex])
		if R.errors[R.runIndex] != nil {
			t.Logf("Executed cmd (%d): ERR: %s", R.runIndex, R.errors[R.runIndex].Error())
		} else {
			t.Logf("Executed cmd (%d)", R.runIndex)
		}
		R.runIndex++
		testSetFlagsDefault()
		cycleCoreForTest()
	}
}

func (R *tRunner) Done() bool { return R.runIndex >= len(R.inputs) }

func runSequential(t *testing.T, runargs ...string) []string {
	output := make([]string, len(runargs))
	for i, a := range runargs {
		s, e := testRoot(a)
		t.Logf("Run %d\nIN: %s\nOUT: %s", i, a, s)
		if e != nil {
			t.Errorf("Run %d Error: %v", i, e)
		}

		output[i] = s
	}
	return output
}
func testClearFlags() { // Run into reinitialization error
	runCmd.ResetFlags()
	specCmd.ResetFlags()
	configCmd.ResetFlags()

	runMakeFlags()
	specMakeFlags()
	configMakeFlags()
}

func testSetFlagsDefault() {
	// add more if causing issues
	*specOps.flags.delete = false
	*specOps.flags.yconfirm = false
	*mainRun.flagSelected = false
	*mainRun.flagAll = false
	*mainRun.fManualRun = false
	*mainRun.fPartialRun = false

}

func cycleCoreForTest() {
	dscore.EndEncode()
	dscore.TempData().Modified = false
	// Any Other Resets Needed?
	configLoadInit()
	dscore.InitTempData()
}

// idk
func testClearFlag(cmd *cobra.Command, flag string) {
	cmd.Flag(flag).Value.Set("")
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

func TestRunnerExecution(t *testing.T) {
	execArgs := []string{"spec fortest", "sel",
		"src d:/coding/examplefiles/OUTPUT", "tgt d:/coding/exampleFiles/big",
		"src d:/coding/examplefiles/OUTPUT", "tgt d:/coding/exampleFiles/OUTPUT", "spec",
	} //"spec fortest --delete -y"}
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

func TestTestReset(t *testing.T) {
	ins := []string{"spec deltest", "spec deltest --delete -y", "spec harblongino"}
	run := testCmdRunner(ins)
	run.ExecuteLog(t)
	t.Logf("%v", *run)

}

func TestFeatureset(t *testing.T) {
	/* When will confirmation be required:
	- [1] shouldn't but does - fix
	- [3]
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
		// BUG: Config Change; Prefs.setOpt() assigns to nil map
		// [7] set cfg for test-audiodirs
		"cfg dry true makealldirs true",
		// [8] add tgt to test-audiodirs
		"tgt d:/coding/exampleFiles/OUTPUT/audio-structure",
		// [9] list (check output correct)
		"list",
		//cleanup
		"spec test-img test-audio test-audiodirs --delete -y",
	}
	//TODO: Fix Cobra/Command/Flag variables not being cleared/initialized each run
	run := testCmdRunner(in)
	run.verbose = true
	run.ExecuteLog(t)
	t.Logf("%v", *run)

	//testRunList(in, t) // Same Bug

	// for i, o := range run.outputs {
	// 	if e := run.errors[i]; e != nil {
	// 		t.Errorf("Run %d Error: %v", i, e)
	// 	}
	// 	t.Logf(`Run %d
	// INPUT: %s
	// OUTPUT: %s`, i, in[i], o)
	// }

}

/*
good
IN: spec test-sound test-svg | OUT: new specs made:\n ***test-sound \n test-svg
------------------------------
good
IN: 'spec test-img --src='...\svg-x-circle' --tgt=.../OUTPUT/images -y'
OUT: spec test-img created and selected
------------------------------
IN: 'src d:/coding/exampleFiles/imagesets/svg_png'
OUT: -- add source(s) --
spec test-img:
        d:/coding/exampleFiles/imagesets/svg_png: true
------------------------------
IN: 'spec test-sound test-svg --delete -y'
OUT: Deleting Specs...
        deleted spec 'test-sound'
        deleted spec 'test-svg'
------------------------------
IN: 'spec test-audio test-audiodirs --src=d:/coding/exampleFiles/audio -y'
OUT: Delete canceled/failed: No Specs found for args.
------------------------------
IN: 'tgt d:/coding/exampleFiles/OUTPUT/audio --ignore=*.mp3'
OUT: 1 ignore patterns added to target d:\coding\exampleFiles\OUTPUT\audio
------------------------------
IN: 'sel iodir'
OUT: error while selecting (No Match Found)
------------------------------
IN: 'cfg dry makealldirs'
OUT: no config options could be made from argscheck cfg --help for argument info
------------------------------
IN: 'tgt d:/coding/exampleFiles/OUTPUT/audio-structure'
OUT: 1 ignore patterns added to target d:\coding\exampleFiles\OUTPUT\audio-structure
------------------------------
IN: 'list'
OUT: User Specs:
*** wez, [3 sources][2 targets] ***
tex, [1 src: D:\Gamedev\textures][1 tgt: D:\Gamedev\pixel-textures](overrides on)
fortest, [1 src: d:\coding\examplefiles\OUTPUT][2 targets]
test-img, [3 sources][1 tgt: d:\coding\exampleFiles\OUTPUT\images]
------------------------------
IN: 'spec test-img test-audio test-audiodirs --delete -y'
OUT: Deleting Specs...
        deleted spec 'test-img'
*/

/* Output Text:
//TODO: commands are running multiple times. fix

Run 0: spec test-sound test-svg
	OUTPUT: Selected test-sound. //WEIRD BUT I THINK FINE
Run 1: spec test-img --src='d:/coding/exampleFiles/imagesets/svg-x-circle,d:/coding/exampleFiles/imagesets/svg_circle' --tgt=d:/coding/exampleFiles/OUTPUT/images -y
	OUTPUT: Selected test-img. AGAIN WEIRD
 Run 2: src d:/coding/exampleFiles/imagesets/svg_png
OUTPUT: -- add source(s)
spec test-img:
	d:/coding/exampleFiles/imagesets/svg_png: false (path exists as source or target in spec)
 Run 3: spec test-sound test-svg --delete -y
	OUTPUT: Deleting Specs... // GOOD
	deleted spec 'test-sound'
	deleted spec 'test-svg'
 Run 4: spec test-audio test-audiodirs --src=d:/coding/exampleFiles/audio -y
	OUTPUT: Delete canceled/failed: No Specs found for args. //NOTE: WTF
 Run 5: tgt d:/coding/exampleFiles/OUTPUT/audio --ignore=*.mp3
	OUTPUT: -- add target(s) --
spec test-imagesets-other:
	d:/coding/exampleFiles/OUTPUT/audio: false (path exists as source or target in spec)
cmd-root_test.go:255: Run 6
	INPUT: sel iodir
	OUTPUT: error while selecting (No Match Found)
Run 7:  cfg dry makealldirs //BUG: Check this
	OUTPUT: no config options could be made from argscheck cfg --help for argument info
Run 8: tgt d:/coding/exampleFiles/OUTPUT/audio-structure
	OUTPUT: -- add target(s) --
spec test-imagesets-other:
	d:/coding/exampleFiles/OUTPUT/audio-structure: false (path exists as source or target in spec)
Run 9: list
	OUTPUT: User Specs:
        *** test-imagesets-other, [3 sources][1 tgt: d:\coding\exampleFiles\OUTPUT\ImageSets] ***
        wez, [1 src: C:\Users\derek\.config\wezterm][0 targets]
        test-img, [2 sources][0 targets]
    cmd-root_test.go:255: Run 10
                INPUT: spec test-img test-audio test-audiodirs --delete -y
                OUTPUT: Deleting Specs...
                deleted spec 'test-img'

*/
