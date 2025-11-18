package cmd

import "testing"

func TestSpecMulti(t *testing.T) {
	run := testCmdRunner([]string{`spec test-img --src='d:/coding/exampleFiles/imagesets/svg-x-circle, d:/coding/exampleFiles/imagesets/svg_circle' -y`})
	run.ExecuteLog(t)
	t.Logf("Input:\n%s", run.inputs[0])
	t.Logf("Output:\n%s", run.outputs[0])
}
