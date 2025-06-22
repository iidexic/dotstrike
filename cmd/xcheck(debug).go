/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	pops "iidexic.dotstrike/pathops"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check called")
		if len(args) > 0 {
			for i, arg := range args {
				cmd.Print("[", i, "] ")
				print(pops.CheckPath(arg))

			}

		}
		//NOTE: this setup might cause some weirdness.
		// -- GetTempGlobals is only supposed to be used when an edit is occurring
		// -- probably best to initialize TempGlobals in a different way
		if *showtempg && dscore.IsTempData() {
			cmd.Printf("%+v", dscore.GetTempGlobals())
		} else if *showtempg {
			cmd.Println("no pending changes (temp is empty)")
		}
	},
}

var (
	walkb, showtempg *bool
)

func init() {
	rootCmd.AddCommand(checkCmd)
	walkb = checkCmd.Flags().BoolP("walk", "w", false, "walk dir")
	showtempg = checkCmd.Flags().Bool("temp", false, "show contents of temporary storage struct for changes to user data")
}
