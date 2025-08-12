package cmd

import (
	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

func warn(text string, tomlWarning bool, cmd *cobra.Command) {
	if tomlWarning {
		cpath := dscore.ConfigTomlPath()
		cmd.Printf(`WARINING:
------------------
%s
Possible error in config file:
(%s)
------------------`, text, cpath)
	} else {
		cmd.Printf("WARINING:\n------------------\n%s\n------------------", text)
	}
}

func warnNilSelectedSpec(cmd *cobra.Command) {
	warn("Selected Spec = Nil.", true, cmd)
}
