/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

var srcF = compFlags{}

/*
Flags that actually apply:
-a --all: print details, --spec
local:
-y --yes, --delete, --alias, --ignore (not implemented)
To be added:
--verbose
--global(?)
*/

var src cmdWrapper

//NOTE: src is the main/master command, tgt is nearly a duplicate. Vars and functions required to run the command are stored here. In the future, they may be combined.
//NOTE: Having a central data structure and building out the command action seems better than this switch

// srcCmd represents the src command
var srcCmd = &cobra.Command{
	Use:   "src path ...",
	Short: "manage sources of a given spec",
	Long:  `add, modify, and delete source components of selected/specified spec`,

	Run: func(cmd *cobra.Command, args []string) {
		affectedSpecs := getSpecs(cmd, true)
		detail, oneOrMoreExist := detailsIfArgsExist(args, affectedSpecs)
		numargs, numspecs := len(args), len(affectedSpecs)
		switch {
		case numargs > 0 && !oneOrMoreExist:
			if oneSpecOrUserConfirm("Adding source to Multiple specs", affectedSpecs) {
				for i := range affectedSpecs {
					added := affectedSpecs[i].CheckAddMultiplePaths(args, true)
					cmd.Printf("Spec %s:\n", affectedSpecs[i].Alias)
					printNumberedListFiltered(cmd, args, added)

				}
			}
		case *srcF.delete && numspecs > 0:
			if oneSpecOrUserConfirm("Deletion with multiple sources or specs", affectedSpecs) && len(args) > 0 {
				for i := range affectedSpecs {
					runDelete(affectedSpecs[i], args, true)
				}
			} else if (*pFlags.all || numspecs == 1) &&
				checkConfirm(fmt.Sprintf("Deletion of ALL sources for %d specs", numspecs), srcF.y) {
				for i := range affectedSpecs {
					affectedSpecs[i].WipeComponentList(true)
				}
			}
		case *pFlags.all && len(args) == 0:
			cmd.Print(detailAllComponentFrom(affectedSpecs, true))
		case len(args) > 0:
			cmd.Print(detail)
		case len(args) == 0 && pFlags.countFlags == 0:
			cmd.Help()
		}

	},
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

func oneSpecOrUserConfirm(requestText string, specs []*dscore.Spec) bool {
	ls := len(specs)
	return ls == 1 || (ls > 1 && checkConfirm(requestText, srcF.y))
}

// !!TODO:(HIGHEST:FIX) DECIDE + IMPLEMENT IF GETSPECS INHERENTLY INCLUDES SELECTED & IF GETSPECS PROCESSES NOSELECT
// getSpecs compiles the list of specs from args of the spec flag
func getSpecs(cmd *cobra.Command, includeSelected bool) []*dscore.Spec {
	specs := []*dscore.Spec{}
	if len(*srcF.spec) > 0 {
		for _, a := range *srcF.spec {
			if s := dscore.TempData().GetSpec(a); s != nil {
				specs = append(specs, s)
			} else {
				cmd.Printf("spec %s not found\n", a)
			}
		}
	}
	if includeSelected {
		specs = append(specs, dscore.TempData().SelectedSpec())
	}
	return specs
}

func detailsIfArgsExist(args []string, specs []*dscore.Spec) (string, bool) {
	if len(args) == 0 {
		return "", false
	}
	//TODO: un-bad this function

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

type compFlags struct {
	alias  *string
	ignore *[]string
	delete *bool
	y      *bool
	spec   *[]string
}

func init() {

	rootCmd.AddCommand(srcCmd)
	srcF.ignore = srcCmd.Flags().StringArray("ignore", nil, "ignore")
	srcF.alias = srcCmd.Flags().String("alias", "", "set alias")
	srcF.delete = srcCmd.Flags().Bool("delete", false, "delete")
	srcF.y = srcCmd.Flags().BoolP("yes", "y", false, "Auto-confirm on prompt")
	srcF.spec = srcCmd.Flags().StringSlice("spec", []string{}, `--spec="alias1, alias2"  to target specs provided`)
}
