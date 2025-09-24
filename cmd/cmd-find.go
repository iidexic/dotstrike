/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find [text] ...",
	Short: "Search existing data for components with matching alias.",
	Long: `Search cfgs, sources, and targets by alias pattern(s) provided as args.

by default, search will return matches from any component (cfg, src, tgt).
find can be restricted to certain contexts or component types using args and flags


flags:
	'--spec -c' search for specs
	'--src -s' search for sources
	'--tgt -t' search for targets

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
		find.cmd = cmd
		find.args = args

		find.runSearch(find.whereLook())
		if len(args) == 0 {
			cmd.Print(cmd.Short)
		}
		fo := cmd.Flags().Lookup("sources")
		if fo != nil {

		}
	},
}

func (f findData) whereLook() (bool, bool, bool) {
	var lookSpec, lookSrc, lookTgt bool
	fcount := 0
	if *f.spec {
		lookSpec = true
		fcount++
	}
	if *f.source {
		lookSrc = true
		fcount++
	}
	if *f.source {
		lookTgt = true
		fcount++
	}
	if fcount == 0 {
		lookSpec, lookSrc, lookTgt = true, true, true
	}
	return lookSpec, lookSrc, lookTgt

}

func (f findData) runSearch(onspec, onsource, ontarget bool) map[string]string {
	temp := dscore.TempData()
	lookup := make(map[string]string, temp.CountComponents()) //test with and without this CountComponents?
	for _, a := range f.args {
		if _, exist := lookup[a]; !exist {
			lookup[a] = fmt.Sprintf("[Find '%s']\n", a)

			if onspec {
				result, mspec := dscore.FindSpec(a)
				switch result {
				case 1:
					lookup[a] += fmt.Sprintf("spec: %s", mspec)
				case -1:
					if *f.fuzzy {
						lookup[a] += fmt.Sprintf("spec (partial match): %s", mspec)
						break
					}
					fallthrough
				case 0:
					lookup[a] += "spec: no match"
				default:
					lookup[a] += "SOMETHING is wrong with find runsearch() or dscore.FindSpec() :)"
				}
			}
		}
	}
	if onsource {

	}
	if ontarget {

	}
	return lookup
}

type findData struct {
	spec, source, target, fuzzy *bool
	cmd                         *cobra.Command
	args                        []string
}

var find findData

func init() {
	rootCmd.AddCommand(findCmd)
	find.spec = findCmd.Flags().Bool("spec", false, "limit search to specs")
	find.source = findCmd.Flags().Bool("src", false, "limit search to sources")
	find.target = findCmd.Flags().Bool("tgt", false, "limit search to targets")
	find.fuzzy = findCmd.Flags().BoolP("ff", "f", false, "fuzzy/inexact find. Currently limited")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
