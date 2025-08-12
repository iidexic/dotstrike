/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Search existing data for components with matching alias.",
	Long: `Search cfgs, sources, and targets by alias pattern(s) provided as args.

by default, search will return matches from any component (cfg, src, tgt).
find can be restricted to certain contexts or component types using args and flags


flags:
	'--cfgs -c' search for cfgs
	'--sources --src -s' search for sources
	'--targets --tgt -t' search for targets

use:
search for all components with alias containing "vscode" or "vim"
	'>ds find vscode vim'

search for sources containing ".config" and targets containing "backup"
	'>ds find -s .config -t backup'

search for all sources of cfg "zsh" containing "data"
	'>ds find zsh -s data
`,
	Run: func(cmd *cobra.Command, args []string) {
		// default behavior
		/* if len(args) == 0{
			cmd.Print(cmd.Short)
		} */
		fo := cmd.Flags().Lookup("sources")
		if fo != nil {

		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
