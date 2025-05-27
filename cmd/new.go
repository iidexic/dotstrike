/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newSCmd represents the new command
var newSCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new group/config",
	Long: `e a new group.
Default args:
	[1] group name
	[2] path

	`,
	Run: actionGroupAdd,
}

func init() {
	groupCmd.AddCommand(newSCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func actionGroupAdd(cmd *cobra.Command, args []string) {
	fmt.Println("group-new: cmd->", cmd, " args=", args)
}
