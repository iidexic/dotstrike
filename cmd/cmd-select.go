/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "sel [pattern]",
	Short: "Selects the first spec whose name/alias contains [pattern]",
	Long: `Select('sel') is used to check or modify the currently selected Spec.

RUNNING:
	Select will check spec aliases against [pattern] arg.
	The first spec with [pattern] as a substring of its alias is selected.
	i.e. '> ds sel a' selects the first spec with 'a' in its name.

By default, all commands will operate on the current selection.
(selected spec will be marked with asterisks when displayed) 
All modifying commands take a '--spec' flag that takes priority over the selected spec.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		temp := dscore.TempData()
		if len(args) == 0 {
			selected := temp.SelectedSpec()
			cmd.Printf("Currently Selected: %s", selected.Alias)
			//cmd.Println("Add partial string to change selection")
			return nil
		}
		for _, arg := range args {
			newsel, e := temp.SelectFirstMatch(arg)
			if e != nil {
				cmd.Printf("error while selecting (%s)\n", e.Error())
			}
			if newsel != "" {
				cmd.Printf("Spec '%s' selected.", newsel)
				return nil
			}
		}
		_ = temp
		return nil
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
