/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

// specCmd represents the cfg command
var specCmd = &cobra.Command{
	Use: "spec alias [sourcePath ... targetPath]",
	Example: `	Create new spec 'term' with 2 sources and a target:
		ds spec term c:\usr\term ~\term c:\globalstore\
			*note: when given 2+ paths, the final path will be set as the spec target.
			 all other paths will be set as sources.
		ds spec term vim`,
	Short: "make a new spec or view existing spec details",
	Long:  `make one new spec, or view details of one or more existing specs.`,
	Run:   specRun,
}

func init() {
	rootCmd.AddCommand(specCmd)
	flagDataSpec = specFlags{
		modify:   *specCmd.Flags().BoolP("modify component", "m", false, "modify"),
		yconfirm: *specCmd.Flags().BoolP("autoconfirm user y/n prompts", "y", false, "yes"),
	}
}

type specFlags struct {
	modify, yconfirm bool
}

var ErrSpecNotMade = errors.New("No spec created; received nil pointer")

var flagDataSpec specFlags
var specOps = specOpData{flags: &flagDataSpec}

func specRun(cmd *cobra.Command, args []string) {

	specOps.args = args
	specOps.argcount = len(args)
	specOps.cmd = cmd
	specOps.argExists = make([]bool, len(args))

	notFound := specOps.populateExisting(args)

	switch {
	case len(args) == 0:
		cmd.Help()
	case len(notFound) == len(args):
		err := specOps.specNew()
		if err != nil {
			cmd.PrintErr(err)
		}
	case len(specOps.existingSpecs) > 0:
		specOps.outputExistingSpecDetails()
	}

}

func (op *specOpData) specNew() error {
	upargs := make([]string, op.argcount)
	copy(upargs, op.args)
	tempdat := dscore.GetTempData()

	var spec *dscore.Spec
	var err error
	if op.argcount > 1 {
		spec, err = tempdat.NewSpec(op.args[0], op.args[1:]...)
	} else {
		spec, err = tempdat.NewSpec(op.args[0])
	}
	if err != nil || spec == nil {
		return fmt.Errorf("error specNew(): %w, from NewSpec: %w", ErrSpecNotMade, err)
	}

	return nil
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
