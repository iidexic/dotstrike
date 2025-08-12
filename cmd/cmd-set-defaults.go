/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tgtCmd represents the tgt command
var defCmd = &cobra.Command{
	Use:   "set-config",
	Short: "Persistently modify default dotstrike configuration",
	Long: `Persistently modify default configuration, which defines how dotstrike operates.
	Add command help here`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("editing  main dotstrike config\nARGS:")
		for i, a := range args {
			cmd.Println(i, ") ", a)
		}
	},
}

func init() {
	specCmd.AddCommand(defCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tgtCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tgtCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
