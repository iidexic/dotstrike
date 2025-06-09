/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

type cfgFlags struct {
}

// cfgCmd represents the cfg command
var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to uickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		lflag := cmd.Flags()
		_ = lflag
	},
}

func init() {
	rootCmd.AddCommand(cfgCmd)
	// PERSISTENT FLAGS:
	//

}
