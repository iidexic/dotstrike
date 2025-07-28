/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

// specCmd represents the cfg command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "make a new spec or view existing spec details",
	Long: `make one new spec, or view details of one or more existing specs.

Usage:
	Make a new spec
	> 'ds spec <alias> [<source paths> <target path>] [options]'
	

	> 'ds spec <alias> [<additional aliases>] [options]'`,
	Run: specRun,
}

type specFlags struct {
	modify, yconfirm bool
}

var flagDataSpec specFlags
var specOps = specOpData{flags: &flagDataSpec}

func specRun(cmd *cobra.Command, args []string) {
	// if not a spec: make new spec
	// if spec exists, show info
	// opSpec := specOpData{
	// 	cmd: cmd, flags: flagDataSpec,
	// 	args: args, argcount: len(args),
	// 	argExists: make([]bool, len(args)),
	// }
	specOps.args = args
	specOps.argcount = len(args)
	specOps.cmd = cmd
	specOps.argExists = make([]bool, len(args))

	notFound := specOps.populateExisting(args) // may not be necessary to return notExists
	//if same number of notFound and args
	switch {
	case len(args) == 0:
		cmd.Help()
	case len(notFound) == len(args):
	case len(specOps.existingSpecs) > 0:
		specOps.outputExistingSpecDetails()
	}
	_ = notFound

}

type specOpData struct {
	flags         *specFlags
	cmd           *cobra.Command
	args          []string
	argcount      int
	argExists     []bool
	existingSpecs []*dscore.Spec
}

func (op *specOpData) populateExisting(args []string) []string {
	notFound := make([]string, 0, len(args))
	for i := range args {
		_ = i
		if spec := dscore.GetSpec(args[i]); spec != nil {
			op.existingSpecs = append(op.existingSpecs, spec)
			op.argExists[i] = true
		} else {
			notFound = append(notFound, args[i])
		}
	}
	return notFound
}
func (op *specOpData) pprintExisting() {
	for i, spec := range op.existingSpecs {
		op.cmd.Print(i, ".   ", spec.Detail())

	}
}

// outputExistingSpecDetails called when have existingSpecs; determines and executes next process steps
func (op *specOpData) outputExistingSpecDetails() {
	if numExist := len(op.existingSpecs); op.argcount == 1 || numExist == op.argcount {
		//all args are existing spec names
	} else if op.argcount > 1 || numExist < op.argcount {
		op.pprintExisting()
	}
}

func countModFlags() int {
	nlocal := 0
	// if flagDataSpec.modify{
	// 	nlocal++
	// }
	return nlocal + pData.countFlags
}

func init() {
	rootCmd.AddCommand(specCmd)
	flagDataSpec = specFlags{
		modify:   *specCmd.Flags().BoolP("modify component", "m", false, "modify"),
		yconfirm: *specCmd.Flags().BoolP("autoconfirm user y/n prompts", "y", false, "yes"),
	}

	// PERSISTENT FLAGS:
	//

}
