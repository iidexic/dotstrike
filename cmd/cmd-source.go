/*
Copyright © 2025 derek :)
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

type componentCmd struct {
	*cmdData
	alias        *string
	ignore       *[]string
	delete       *bool
	y            *bool
	spec         *[]string
	selectedSpec *bool
	isSource     bool
}

// srcCmd represents the src command
var srcCmd = &cobra.Command{
	Use:   "src path ...",
	Short: "manage sources of a given spec",
	Long:  `add, modify, and delete source components of selected/specified spec`,
	Run:   sourceRun,
}
var src = componentCmd{}

func sourceRun(cmd *cobra.Command, args []string) {
	src.cmdData = newCmdData(cmd, args)
	specFlagArgs := *src.spec

	var ns = 0

	if len(specFlagArgs) > 0 {
		ns = src.getSpecs(false, specFlagArgs...)
	} else {
		ns = src.getSpecs(false)
	}

	if len(args) > 0 {
		for _, spec := range src.specs {
			src.components = append(src.components,
				spec.GetMatchingComponents(src.args, src.isSource)...)
		}

	}

	if ns > 0 {
		e := runComponent(&src)
		if e != nil {
			cmd.Print(e.Error())
		}
	} else if ns == 0 {
		cmd.Print("no specs found for entered arguments!")
	}

}

func runComponent(cmp *componentCmd) error {

	nargs, ncomp := len(cmp.args), len(cmp.components)
	if nargs > 0 && ncomp == 0 { // Make new args
	} else if nargs > 0 && ncomp < nargs { // Uncertain/Mixed Qty
		if *cmp.delete {
			return error(fmt.Errorf("Found no matching sources to delete"))
		}
		cmp.runMixedQty()
	} else if ncomp > 0 {
		switch {
		case *cmp.delete && len(*cmp.ignore) > 0: //Delete Ignores
			if len(cmp.components) == 0 {

			}

		case *cmp.delete: //Delete Sources

		case len(*cmp.ignore) > 0: // Add Ignores
			return cmp.addIgnores()
		}

	}

	return nil
}

func (C *componentCmd) runMixedQty() error {
	switch {
	case len(*C.ignore) > 0 && *C.delete:
		for i, comp := range C.components {
			if len(comp.Ignores) > 0 {
				_ = i
			}
		}
	case len(*C.ignore) > 0:

	}

	return nil
}

func (C *componentCmd) deleteArgs() int {

	return 0
}

func (C *componentCmd) addIgnores() error {

	return nil
}

func (C *componentCmd) deleteIgnores() error {

	return nil
}

func runDelete(spec *dscore.Spec, args []string, isSource bool) {
	for _, arg := range args {
		result := spec.DeleteIfChild(arg, isSource, false)
		if result == 0 {
			src.Printf("spec %s: None deleted", spec.Alias)
		} else {
			src.Printf("spec %s: %d deleted", spec.Alias, result)
		}
	}
}

var (
	msgAddSource       = "Add source(s) to multiple specs"
	msgAddIgnores      = "Add Ignores to multiple sources"
	fmsgDelSourceCount = "Delete %d Source(s)"
	msgDelIgnores      = "Delete Ignore pattern(s) in multiple sources"
)

func oneSpecOrUserConfirm(requestText string, specs []*dscore.Spec) bool {
	ls := len(specs)
	return ls == 1 || (ls > 1 && checkConfirm(requestText, src.y))
}

// !!TODO:(HIGHEST:FIX) DECIDE + IMPLEMENT IF GETSPECS INHERENTLY INCLUDES SELECTED & IF GETSPECS PROCESSES NOSELECT
// getSpecs compiles the list of specs from args of the spec flag
func getSpecs(cmd *cobra.Command, includeSelected bool) []*dscore.Spec {
	specs := []*dscore.Spec{}
	if includeSelected {
		specs = append(specs, dscore.TempData().SelectedSpec())
	}
	if len(*src.spec) > 0 {
		for _, a := range *src.spec {
			if s := dscore.TempData().GetSpec(a); s != nil {
				specs = append(specs, s)
			} else {
				cmd.Printf("spec %s not found\n", a)
			}
		}
	}
	return specs
}

func detailsIfArgsExist(args []string, specs []*dscore.Spec) (string, bool) {
	if len(args) == 0 {
		return "", false
	}
	//TODO: un-bad this function
	// This is completely unreadable

	// Functionality:
	// 1. collects details of the specs passed, prints them directly
	// 2. checks if args include identifier of an existing component within given specs
	hasExisting := false
	details := make([]string, 0, len(specs)+len(args))
	for _, spec := range specs {
		if existing := spec.GetExistingChildren(args); len(existing) > 0 {
			hasExisting = true
			details = append(details, "In spec "+spec.Alias+":")
			srctxt, tgttxt := "", ""
			for i, component := range existing {
				if component.IsSource() {
					srctxt += fmt.Sprintf(" %d. %s:%s\n", i, component.Alias, component.Path)
				} else {
					tgttxt += fmt.Sprintf(" %d. %s:%s", i, component.Alias, component.Path)
				}
			}
			details = append(details, "[sources]", srctxt, "\n[targets]", tgttxt)
		}

	}
	return strings.Join(details, "\n"), hasExisting
}

func detailAllComponentFrom(specs []*dscore.Spec, isSource bool) string {

	detail := make([]string, 0, len(specs)*2) //arbitrary
	for _, spec := range specs {
		detail = append(detail, fmt.Sprintf("Spec: %s", spec.Alias))
		if isSource {
			if len(spec.Sources) == 0 {
				detail = append(detail, " [no sources]")

			} else {
				for i := range spec.Sources {
					detail = append(detail, spec.Sources[i].Detail())
				}
			}
		} else {

			if len(spec.Targets) == 0 {
				detail = append(detail, " [no targets]")

			} else {
				for i := range spec.Targets {
					detail = append(detail, spec.Targets[i].Detail())
				}
			}
		}
	}
	return strings.Join(detail, "\n")

}

func init() {

	rootCmd.AddCommand(srcCmd)
	src.ignore = srcCmd.Flags().StringArray("ignore", nil, "ignore")
	src.alias = srcCmd.Flags().String("alias", "", "set alias")
	src.delete = srcCmd.Flags().Bool("delete", false, "delete")
	src.y = srcCmd.Flags().BoolP("yes", "y", false, "Auto-confirm on prompt")
	src.spec = srcCmd.Flags().StringSlice("spec", []string{}, `--spec="alias1, alias2"  to target specs provided`)
	src.selectedSpec = srcCmd.Flags().BoolP("selected", "s", true, "--useSelected=false to disable operating on selected spec")
}
