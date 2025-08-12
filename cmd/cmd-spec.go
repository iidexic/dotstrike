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

var temp dscore.Temp

//TODO: (mid-doc) Update/Correct Long detail to match current function

// specCmd represents the cfg command
var specCmd = &cobra.Command{
	Use: "spec alias [sourcePath ... targetPath]",
	Example: `	Create new spec 'term' with 2 sources and a target:
		ds spec term c:\usr\term ~\term c:\globalstore\
			*note: when given 2+ paths, the final path will be set as the spec target.
			 all other paths will be set as sources.
		ds spec term vim`,
	Short: "make a new spec by providing a unique alias and paths(optional)",
	Long: `Make one new spec using the first argument as the alias. 
	All specs must have a unique alias; creation will fail if alias is not unique.

Paths can be provided as arguments after the alias argument. 
	If one path is provided: it will be added as a child source to the new spec.
	If two paths are provided: the first path will be added as a source, the second will be added as a target.
	if n > 2 paths are passed: paths 1 to (n-1) will be added as sources, and the final path will be added as a target.`,

	Run: specRun,
}

func init() {
	rootCmd.AddCommand(specCmd)

	flagDataSpec = specFlags{
		delete:   specCmd.Flags().Bool("delete", false, "delete spec"),
		yconfirm: specCmd.Flags().BoolP("autoconfirm user y/n prompts", "y", false, "yes"),
		alias:    specCmd.Flags().String("set-alias", "", "set-alias ALIAS"),
	}
}

type specFlags struct {
	yconfirm, delete *bool
	alias            *string
}

var ErrSpecNotMade = errors.New("No spec created; received nil pointer")

var flagDataSpec specFlags
var specOps = specOpData{flags: &flagDataSpec}

// TODO: source/target flags
// TODO: Use SelectedSpec for 0-arg edits
func specRun(cmd *cobra.Command, args []string) {
	temp = dscore.TempData()
	specOps.args = args
	specOps.argcount = len(args)
	specOps.cmd = cmd
	specOps.argExists = make([]bool, len(args))

	notFound := specOps.populateExisting(args)
	switch {

	case len(args) == 0:
		act := specOps.checkFlagActions()
		if !act {
			specOps.outputSelected()
		}
	case len(notFound) == len(args):
		err := specOps.specNew()
		if err != nil {
			cmd.PrintErr(err)
		}
	case len(specOps.existingSpecs) > 0:
		if *specOps.flags.delete {
			for i := range specOps.existingSpecs {
				deleted := specOps.specDelete(specOps.existingSpecs[i].Alias)
				if !deleted {
					cmd.Printf("delete %s failed.", specOps.existingSpecs[i].Alias)
				}
			}
		} else {
			changed := dscore.TempData().SelectPtr(specOps.existingSpecs[0])
			if changed {
				cmd.Printf("Selected %s.", specOps.existingSpecs[0].Alias)
			}
		}
		// Select
		//specOps.outputExistingSpecDetails()
	}

}
func (op *specOpData) outputSelected() {
	s := dscore.TempData().SelectedSpec()
	if s == nil {
		op.cmd.Printf("error: spec not found. resetting  selection")
		op.cmd.Help()
	} else {
		op.cmd.Print("Selected spec: ", s.Alias, "\n\n")
		speclist := dscore.TempData().Specs
		for i := range speclist {
			op.cmd.Printf("[%d] ", i)
			if &speclist[i] == s {
				op.cmd.Print("***")
			}
			op.cmd.Print(speclist[i].Alias, "\n")
		}
		op.cmd.Print(s.Detail())
	}
}

func (op *specOpData) specNew() error {
	upargs := make([]string, op.argcount)
	copy(upargs, op.args)
	tempdat := dscore.TempData()
	var spec *dscore.Spec
	var err error
	if op.argcount > 1 {
		spec, err = tempdat.NewSpec(op.args[0], op.args[1:]...)
	} else {
		spec, err = tempdat.NewSpec(op.args[0])
	}
	if err != nil || spec == nil {
		return fmt.Errorf("error in op.specNew(): %w, from NewSpec: %w", ErrSpecNotMade, err)
	}

	return nil
}

func (op *specOpData) checkFlagActions() bool {
	if newAlias := *op.flags.alias; newAlias != "" {
		switch {
		case op.argcount == 1 && len(op.existingSpecs) == 1:
			spec := op.existingSpecs[0]
			op.editSpecAlias(spec, newAlias)
		case op.argcount == 0:
			spec := dscore.TempData().SelectedSpec()
			op.editSpecAlias(spec, newAlias)
		case op.argcount > 1:
			op.cmd.Print("Error - cannot change alias of multiple specs!")
		}
		return true
	}
	return false
}

func (op *specOpData) editSpecAlias(spec *dscore.Spec, newAlias string) {
	ogAlias := spec.Alias
	changed := dscore.TempData().ChangeSpecAlias(spec, newAlias)
	if !changed {
		op.cmd.Printf(`Updating spec '%s' to '%s' failed:
Alias not unique (spec '%s' already exists)`, spec.Alias, newAlias, newAlias)
	} else {
		op.cmd.Printf("Spec '%s' updated to '%s'", ogAlias, newAlias)
	}
}

// func (op *specOpData) checkConfirmExecute(fn func(), detail string) {
// 	approve := *op.flags.yconfirm
// 	if !approve {
// 		approve = askConfirmf(detail)
// 	}
// 	if *op.flags.yconfirm || askConfirmf(detail) {
// 		fn()
// 	}
// }

func (op *specOpData) checkConfirm(detail string) bool {
	return *op.flags.yconfirm || askConfirmf(detail)
}

// TODO:(hi-refactor) fix specDelete and/or shift to specPtrDelete as ptr is getting pulled beforehand anyway.
func (op *specOpData) specPtrDelete(spec *dscore.Spec) bool {
	out := false

	if spec != nil {
		if op.checkConfirm("Delete spec " + spec.Alias) {
			out = dscore.TempData().DeleteSpec(spec)
		}
	}
	return out
}

func (op *specOpData) specDelete(alias string) bool {
	sptr := dscore.TempData().GetSpec(alias)
	if sptr == nil {
		op.cmd.PrintErrf("specDelete error: %s", dscore.ErrAliasNotFound.Error())
	} else if op.checkConfirm("Delete spec " + sptr.Alias) {
		dscore.TempData().Modified = true
		return dscore.TempData().DeleteSpec(sptr)
	}
	return false
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
		if spec := dscore.TempData().GetSpec(args[i]); spec != nil {
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
