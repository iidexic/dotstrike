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

//NOTE: Having a central data structure and building out the command action seems way better than this switch
//TODO: Refactor entire source command

// srcCmd represents the src command
var srcCmd = &cobra.Command{
	Use:   "src path ...",
	Short: "manage sources of a given spec",
	Long:  `add, modify, and delete source components of selected/specified spec`,

	Run: func(cmd *cobra.Command, args []string) {
		src = cmdWrapper{Command: cmd, args: args} //WIP/FUTURE IMPLEMENTATION
		affectedSpecs := getSpecs(cmd)
		switch {
		case len(args) > 0 && !detailsIfArgsExist(cmd, args, affectedSpecs):
			if oneSpecOrUserConfirm("Adding source to Multiple specs", affectedSpecs) {
				for i := range affectedSpecs {
					//TODO: standardize modify. in delete its run within the modifying function.
					dscore.TempData().Modify()
					added := affectedSpecs[i].CheckAddMultiplePaths(args, true)
					cmd.Printf("Spec %s:\n", affectedSpecs[i].Alias)
					printNumberedListFiltered(cmd, args, added)

				}
			}
		case *srcF.delete && len(affectedSpecs) > 0:
			if conf := oneSpecOrUserConfirm("Deletion with multiple sources or specs", affectedSpecs); conf && len(args) > 0 {
				for i := range affectedSpecs {
					runDelete(affectedSpecs[i], args, true)
				}
			} else if conf && *pFlags.all {
				for i := range affectedSpecs {
					affectedSpecs[i].WipeComponentList(true)
				}
			}
		case *pFlags.all && len(args) == 0:
			cmd.Print(detailAllComponentFrom(affectedSpecs, true))
		case len(args) == 0 && pFlags.countFlags == 0:
			cmd.Help()
		}

		if *pFlags.all {
		}

	},
}

func runDelete(spec *dscore.Spec, args []string, isSource bool) {

	if isSource {
		for _, arg := range args {
			//test
			src.Printf("trying to delete %s in spec %s\n", arg, spec.Alias)
			src.Printf("count sources = %d\n", len(spec.Sources))
			src.Print(spec.Sources)
			//end test
			result := spec.DeleteIfChild(arg, true)
			src.Printf("\nDeleted? -> %t\n", result)
			src.Printf("count sources = %d\n", len(spec.Sources))
			src.Print(spec.Sources)
		}
	}
}

func oneSpecOrUserConfirm(requestText string, specs []*dscore.Spec) bool {
	ls := len(specs)
	return ls == 1 || (ls > 1 && checkConfirm(requestText, srcF.y))
}

func getSpecs(cmd *cobra.Command) []*dscore.Spec {
	specs := []*dscore.Spec{}
	if pFlags.bspec {
		for _, a := range *pFlags.spec {
			if s := dscore.TempData().GetSpec(a); s != nil {
				specs = append(specs, s)
			} else {
				cmd.Printf("spec %s not found\n", a)
			}
		}
	} else {
		specs = append(specs, dscore.TempData().SelectedSpec())
	}
	return specs
}

func detailsIfArgsExist(cmd *cobra.Command, args []string, specs []*dscore.Spec) bool {
	if len(args) == 0 {
		return true
	}
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
	if hasExisting {
		cmd.Print(strings.Join(details, "\n"))
	}
	return hasExisting
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
}

func init() {

	rootCmd.AddCommand(srcCmd)
	srcF.ignore = srcCmd.Flags().StringArray("ignore", nil, "ignore")
	srcF.alias = srcCmd.Flags().String("alias", "", "set alias")
	srcF.delete = srcCmd.Flags().Bool("delete", false, "delete")
	srcF.y = srcCmd.Flags().BoolP("yes", "y", false, "Auto-confirm on prompt")
}
