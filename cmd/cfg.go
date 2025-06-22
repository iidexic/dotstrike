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
		Subcommands:
		- src
		- tgt
		*/
		//Persistent Flags Short: v, g(?), c, s, t
		//1+configs: if mix of non/existent, fail if context can't be determined

		/* find := dscore.Operation{
			Get: dscore.Lookup{
				GetCfg: len(*pf.cfg) > 0,
				GetSrc: len(*pf.src) > 0,
				GetTgt: len(*pf.tgt) > 0},
		find.ProcessFind()
		} */
		if len(args) > 0 {
			found := dscore.SelectCfg(args[0])
			if !found {
				cmd.Printf("cfg %s not found", args[0])
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
