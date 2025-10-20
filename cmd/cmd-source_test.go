package cmd

import "testing"

func TestSrcCmdBasicAll(t *testing.T) {
	ins := []string{"spec srcTest", "src C:/secret/bringo", "src"}
	_, e := testRunSequence(ins, t)
	if e != nil {
		t.Errorf("Failures during run sequence")
	}

}
