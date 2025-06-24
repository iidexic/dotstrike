/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

type cfgFlags struct {
}

// cfgCmd represents the cfg command
var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "short descr",
	Long:  `long descr`,
	Run: func(cmd *cobra.Command, args []string) {
		/* ╭───────────────────────── CFG command logic ─────────────────────────╮
		Subcommands: - src - tgt
		//Persistent Flags Short: v,  c, s, t, g(remove),
		} */
		if len(args) > 0 {
			found := dscore.SelectCfg(args[0])
			if !found {
				cmd.Printf("cfg %s not found", args[0])
			} else {

			}

		}
	},
}

func init() {
	rootCmd.AddCommand(cfgCmd)
	cfgCmd.Flags().BoolP("modify component", "m", false, "modify")
	cfgCmd.Flags().BoolP("apply to all found", "a", false, "all")

	// PERSISTENT FLAGS:
	//

}
