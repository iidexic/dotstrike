/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var tgt = componentCmd{}

// tgtCmd represents the tgt command
var tgtCmd = &cobra.Command{
	Use:   "tgt path ...",
	Short: "manage targets of a given spec",
	Long:  `add, modify, and delete source components of selected/specified spec`,

	Run: func(cmd *cobra.Command, args []string) {
		tgt.cmdData = newCmdData(cmd, args)
		tgt.isSource = false
		specFlagArgs := *tgt.spec
		ns := tgt.getSpecs(false, specFlagArgs...)
		if len(args) > 0 {
			for _, spec := range tgt.specs {
				tgt.components = append(tgt.components,
					spec.GetMatchingComponents(tgt.args, tgt.isSource)...)
			}
		}

		if ns > 0 {
			e := runComponent(&tgt)
			if e != nil {
				cmd.Print(e.Error())
			}
		} else if ns == 0 {
			cmd.Print("no specs found for entered arguments!")
		}
	},
}

func init() {
	rootCmd.AddCommand(tgtCmd)
	makeCmpFlags(tgtCmd, &tgt)
}
