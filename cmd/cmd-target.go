/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tgtF = compFlags{}

var tgt cmdWrapper

//NOTE: Having a central data structure and building out the command action seems way better than this switch
//TODO: Refactor source command

// tgtCmd represents the tgt command
var tgtCmd = &cobra.Command{
	Use:   "tgt path ...",
	Short: "manage targets of a given spec",
	Long:  `add, modify, and delete source components of selected/specified spec`,

	Run: func(cmd *cobra.Command, args []string) {
		tgt = cmdWrapper{Command: cmd, args: args} //WIP/FUTURE IMPLEMENTATION
		affectedSpecs := getSpecs(cmd, true)
		detail, oneOrMoreExist := detailsIfArgsExist(args, affectedSpecs)
		numargs, numspecs := len(args), len(affectedSpecs)
		switch {
		case numargs > 0 && !oneOrMoreExist:
			if oneSpecOrUserConfirm("Adding target to Multiple specs", affectedSpecs) {
				for i := range affectedSpecs {
					added := affectedSpecs[i].CheckAddMultiplePaths(args, false)
					cmd.Printf("Spec %s:\n", affectedSpecs[i].Alias)
					printNumberedListFiltered(cmd, args, added)
				}
			}
		case *tgtF.delete && numspecs > 0:
			if oneSpecOrUserConfirm(
				fmt.Sprintf("Delete on %d specs", numspecs), affectedSpecs) && len(args) > 0 {
				for i := range affectedSpecs {
					runDelete(affectedSpecs[i], args, false)
				}
			} else if (*pFlags.all || numspecs == 1) &&
				checkConfirm(fmt.Sprintf("Delete ALL targets - %d specs", numspecs), tgtF.y) {
				for i := range affectedSpecs {
					affectedSpecs[i].WipeComponentList(false)
				}
			}
		case *pFlags.all && len(args) == 0:
			cmd.Print(detailAllComponentFrom(affectedSpecs, false))

		case len(args) > 0:
			cmd.Print(detail)
		case len(args) == 0 && pFlags.countFlags == 0:
			cmd.Help()
		}

	},
}

func init() {
	rootCmd.AddCommand(tgtCmd)

	tgtF.ignore = tgtCmd.Flags().StringArray("ignore", nil, "ignore")
	tgtF.alias = tgtCmd.Flags().String("alias", "", "set alias")
	tgtF.delete = tgtCmd.Flags().Bool("delete", false, "delete")
	tgtF.y = tgtCmd.Flags().BoolP("yes", "y", false, "Auto-confirm on prompt")
}
