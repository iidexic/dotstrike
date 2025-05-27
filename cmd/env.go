/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envSCmd = &cobra.Command{
	Use:   "env",
	Short: "show env var created",
	Long:  `show full env; if dotstrike var is included, call it out`,
	Run:   envRun,
}

func envRun(cmd *cobra.Command, args []string) {
	print("no env :(")
}

func init() {
	showCmd.AddCommand(envSCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
