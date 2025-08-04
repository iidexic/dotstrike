/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var tgtF = compFlags{}

/*
Flags that actually apply:
-a --all: print details, --spec
local:
-y --yes, --delete, --alias, --ignore (not implemented)
To be added:
--verbose
--global(?)
*/

var tgt cmdWrapper

//NOTE: Having a central data structure and building out the command action seems way better than this switch
//TODO: Refactor entire source command

// tgtCmd represents the tgt command
var tgtCmd = &cobra.Command{
	Use:   "tgt path ...",
	Short: "manage sources of a given spec",
	Long:  `add, modify, and delete source components of selected/specified spec`,

	Run: func(cmd *cobra.Command, args []string) {
		tgt = cmdWrapper{Command: cmd, args: args} //WIP/FUTURE IMPLEMENTATION
		affectedSpecs := getSpecs(cmd)
		switch {
		case len(args) > 0 && !detailsIfArgsExist(cmd, args, affectedSpecs):
			if oneSpecOrUserConfirm("Adding source to Multiple specs", affectedSpecs) {
				for i := range affectedSpecs {
					affectedSpecs[i].CheckAddMultiplePaths(args, false)
				}
			}
		case *tgtF.delete && len(affectedSpecs) > 0:
			if conf := oneSpecOrUserConfirm("Deletion with multiple sources or specs", affectedSpecs); conf && len(args) > 0 {
				for i := range affectedSpecs {
					for _, arg := range args {
						affectedSpecs[i].DeleteIfChild(arg, false)
					}
				}
			} else if conf && *pFlags.all {
				for i := range affectedSpecs {
					affectedSpecs[i].WipeComponentList(false)
				}
			}
		case *pFlags.all && len(args) == 0:
			cmd.Print(detailAllComponentFrom(affectedSpecs, false))
		case len(args) == 0 && pFlags.countFlags == 0:
			cmd.Help()
		}

		if *pFlags.all {
		}

	},
}

func init() {
	rootCmd.AddCommand(tgtCmd)

	tgtF.ignore = tgtCmd.Flags().StringArray("ignore", nil, "ignore")
	tgtF.alias = tgtCmd.Flags().String("alias", "", "set alias")
	tgtF.delete = tgtCmd.Flags().Bool("delete", false, "delete")
	tgtF.y = tgtCmd.Flags().BoolP("yes", "y", false, "Auto-confirm on prompt")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tgtCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tgtCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
