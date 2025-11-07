/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "sel [pattern]",
	Short: "Select a spec by alias, either selects an exact match or the first substring match",
	Long: `Select('sel') is used to check the currently selected spec, or to select a different spec by alias.

RUNNING:
Select will check spec aliases against the first arg.
	An exact match between arg and spec alias will be attempted first.
	If no match is found, the first spec with the arg text in its alias is selected.
	i.e. '> ds sel a' selects the first spec with 'a' in its alias.

By default, all commands will operate on the current selection.
(selected spec will be marked with asterisks when displayed) 
All modifying commands take a '--spec' flag that takes priority over the selected spec.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		temp := dscore.TempData()
		if len(args) == 0 {
			selected := temp.SelectedSpec()
			if selected == nil {
				return fmt.Errorf("No Spec Selected")
			}
			cmd.Printf("Currently Selected: %s", selected.Alias)
			//cmd.Println("Add partial string to change selection")
			return nil
		}
		if temp.Select(args[0]) {
			cmd.Printf("Spec '%s' selected.", args[0])
			return nil
		}
		for _, arg := range args { //why
			newsel, e := temp.SelectFirstMatch(arg)
			if e != nil {
				cmd.Printf("error while selecting (%s)\n", e.Error())
			}
			if newsel != "" {
				cmd.Printf("Selected '%s' (first with '%s')", newsel, arg)
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
