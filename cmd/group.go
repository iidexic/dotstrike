/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g"},
	Short:   "Command used to make and manage groups",
	Long: `
	group provides subcommands for managing groups of appconfigs.
	the group command can be used in the following ways:
		
	> dg group 'groupName':	switch to existing group; 
				if 'groupName' does not exist, user will be prompted
				to confirm if they would like to create 'groupName'
	
	subcommands:
		'new' - 
		'delete' - 
		'clone' -
		'merge' -
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("group called")
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
