/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	"iidexic.dotstrike/uout"
)

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
	Long: ` Create new specs or delete existing specs.

	Create a new spec by passing a name/unique alias:
		'> ds spec term'
	All specs must have a unique alias, as this is their primary identifier.
	
	Behavior:
	- 1 argument: If matches an existing spec, that spec will be selected.
		If no existing spec matches the alias, a new spec will be created.
	- 2 or more arguments: If the first argument matches an existing spec, that spec will be selected.
		If the first argument is unique, new specs will be created for each unique argument.
	`,

	Run: specRun,
}

func init() {
	rootCmd.AddCommand(specCmd)

	specMakeFlags()

}

func specMakeFlags() {

	flagDataSpec = specFlags{
		delete:   specCmd.Flags().BoolP("delete", "d", false, "delete spec"),
		yconfirm: specCmd.Flags().BoolP("autoconfirm user y/n prompts", "y", false, "yes"),
		alias:    specCmd.Flags().String("set-alias", "", "set-alias ALIAS"),
		src: specCmd.Flags().StringSlice("src", make([]string, 0, 2),
			`--src="c:\srcPath1,.\path2"`),
		tgt: specCmd.Flags().StringSlice("tgt", make([]string, 0, 2),
			`--tgt="c:\target\path1,.\tpath2"`),
		ignore:   specCmd.Flags().StringSlice("ignore", make([]string, 0, 2), "--ignore='ptn1,ptn2'"),
		validate: specCmd.Flags().Bool("validate", false, "validate that spec's source paths exist. If not, remove them."),
	}
	specOps.flags = &flagDataSpec
}

type specFlags struct {
	yconfirm, delete *bool
	alias            *string
	src, tgt         *[]string
	ignore           *[]string
	validate         *bool
}

var ErrSpecNotMade = errors.New("No spec created; received nil pointer")

var flagDataSpec specFlags
var specOps = specOpData{flags: &flagDataSpec}

// TODO:(Done?) Use SelectedSpec for 0-arg edits

func specRun(cmd *cobra.Command, args []string) {
	td := dscore.TempData()
	specOps.args = sliceUniques(args)
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
		if *specOps.flags.delete {
			if e := specOps.uncertainDelete(); e != nil {
				cmd.PrintErr(e)
			}
		} else {
			err := specOps.specNew()
			if err != nil {
				cmd.PrintErr(err)
			}
		}
	case len(specOps.existingSpecs) > 0:
		if *specOps.flags.delete { //TODO: Finish Correcting Spec Delete (processDeletion())
			if e := specOps.processDeletion(); e != nil {
				cmd.Printf("Error during deletion: %s", e.Error())
			}
		} else {
			changed := td.SelectPtr(specOps.existingSpecs[0])
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
		op.cmd.Printf("error: spec not found. resetting selection")
		// WARN: I should not have to do this here, figure out why specdelete fails to correct
		dscore.TempData().Modify()
		dscore.ResetSpecSelection()
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
	temp := dscore.TempData()
	if op.reqMultNewWithPaths() {
		ask := checkConfirmF("make %d specs with same source/target?", op.flags.yconfirm, len(op.args))
		if !ask {
			op.cmd.Print("0 specs made")
			return nil
		}
	}
	speclenPre := len(temp.Specs)
	var err error
	specsMade := make([]string, op.argcount)
	for i, v := range op.args {
		s, e := temp.NewSpec(v, *op.flags.src, *op.flags.tgt)
		if e != nil {
			if err == nil {
				err = fmt.Errorf("err making new spec:")
			} else {
				err = fmt.Errorf("%w, %w", err, e)
			}

		}
		specsMade[i] = s.Alias

	}
	if err != nil {
		op.cmd.Print("Warning: Errors making specs")
	}
	if numNewSpecs := len(temp.Specs) - speclenPre; numNewSpecs > 0 {
		selectionChanged := temp.Select(specsMade[0])
		switch {
		case selectionChanged && numNewSpecs == 1:
			op.cmd.Printf("spec %s created and selected\n", specsMade[0])
		case selectionChanged && numNewSpecs > 1:
			op.cmd.Print("new specs made:\n***")
			for _, alias := range specsMade {
				op.cmd.Printf("%s\n", alias)
			}
		case numNewSpecs == 0 && err != nil:
			op.cmd.Print("No new specs. ERRORS:")
			op.cmd.Print(err.Error())

		case numNewSpecs == 0:
			op.cmd.Print("No new specs. ERRORS:")
		}
	}
	return nil
}

func (op *specOpData) validateSpecs() {
	for _, s := range op.existingSpecs {
		badPaths := s.ValidateAndCleanSources()
		op.cmd.Printf("Spec %s: Removed %d bad paths:\n%s", s.Alias, len(badPaths), badPaths)
	}
}

// checks args and flag args, returns true if more than 1 main arg and at least one path flag arg
func (op specOpData) reqMultNewWithPaths() bool {
	nArgs := len(op.args)
	nSrcArgs := len(*op.flags.src)
	nTgtArgs := len(*op.flags.tgt)
	// true if more than 1 main arg and at least one path flag arg
	return nArgs > 1 && (nSrcArgs > 0 || nTgtArgs > 0)
}

func (op *specOpData) checkFlagActions() bool {

	// alias change
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
	if *op.flags.validate {
		op.validateSpecs()
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

func (op *specOpData) runargCheck() {
	rootCmd.InOrStdin()
}

func (op *specOpData) checkConfirm(detail string) bool {
	return checkConfirmF(detail, op.flags.yconfirm)
}

func (op *specOpData) processDeletion() error {
	out := uout.NewOutf("Delete %d Specs. All specs to be deleted:", len(op.existingSpecs))
	out.IndR()
	aliases := make([]string, len(op.existingSpecs))
	for i := range op.existingSpecs {
		s := op.existingSpecs[i]
		aliases[i] = s.Alias
		out.F("%s (%d Sources, %d Targets)", s.Alias, len(s.Sources), len(s.Targets))
	}
	out.WipeOnOutput(true)
	if op.checkConfirm(out.String()) {
		out.IndL().A("Deleting Specs...")
		out.IndR()
		deleted := dscore.TempData().DeleteSpecs(aliases)
		out.IfLN(deleted, "deleted spec '%s'", "failed to delete spec '%s'", aliases)
		op.cmd.Print(out.String())
	}
	return nil
}

func (op *specOpData) uncertainDelete() error {
	aliases := dscore.TempData().SubstringSearchSpecs(op.args)
	out := uout.NewOut("Delete attempted, no exact alias matches found.")
	out.V("Found Specs:")
	out.IndR().LV(aliases)
	out.IndL().V("CONFIRM DELETION OF FOUND SPECS")
	if len(aliases) > 0 && op.checkConfirm(out.String()) {
		out.Clear()
		out.A("Deleting Specs...")
		deleted := dscore.TempData().DeleteSpecs(aliases)
		out.IfLN(deleted, "deleted spec '%s'", "failed to delete spec '%s'", aliases)
		op.cmd.Print(out.String())
	} else {
		op.cmd.Print("Delete canceled/failed: No Specs found for args.")
	}
	return nil
}

// specDelete will perform the required steps to delete an individual spec
// will not check for confirmation
// func (op *specOpData) specDelete(alias string) bool { return dscore.TempData().DeleteSpec(alias) }

type specOpData struct {
	flags         *specFlags
	cmd           *cobra.Command
	args          []string
	argcount      int
	argExists     []bool
	existingSpecs []*dscore.Spec
}

func (op *specOpData) populateExisting(args []string) []string {

	notFound := make([]string, len(args))
	op.existingSpecs = make([]*dscore.Spec, len(args))
	ix, in := 0, 0
	for i, s := range args {
		if spec := dscore.TempData().GetSpec(s); spec != nil {
			op.existingSpecs[ix] = spec
			ix++
			op.argExists[i] = true
		} else {
			notFound[in] = s
			in++
		}
	}
	op.existingSpecs = op.existingSpecs[:ix]
	if *persistentFlags.verbose {
		op.cmd.Printf("found %d specs", ix)
	}
	return notFound[:in]
}
