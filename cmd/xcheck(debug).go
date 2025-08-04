/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	pops "iidexic.dotstrike/pathops"
)

type checkCmdFlags struct {
	show, temp, walk, ask *bool
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// verbose "v", globals "g", cfg, src, tgt "c","s","t"
	Run: func(cmd *cobra.Command, args []string) {
		if *checkf.ask {

			conf := askConfirmf("Is this question true")
			fmt.Printf("function gave: %t\n", conf)
			confright := askConfirmf("Was that correct")
			if confright {
				print("thats good")
			} else {
				print("oh no")
			}
		} else {
			fmt.Println("check called")
			if len(args) > 0 {
				for i, arg := range args {
					cmd.Print("[", i, "] ")
					print(pops.CheckPathDebug(arg))

				}

			}
			if pFlags.countFlags > 0 {

			}
			if td := dscore.TempData(); td != nil && *checkf.temp && td.Modified {
				cmd.Printf("%+v", td)
			} else if *checkf.temp {
				cmd.Println("no pending changes (temp is empty)")
			}
			if *checkf.show {
				exec, e := os.Executable()
				ce(e)
				fmt.Println(exec)
				wd, e := os.Getwd()
				ce(e)
				fmt.Printf("workingDir: %s\n", wd)
				callfrom := pops.CalledFrom()
				fmt.Printf("args[0]- called from: %s\n", callfrom)
			}
		}

	},
}

var checkf = checkCmdFlags{}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkf.walk = checkCmd.Flags().BoolP("walk", "w", false, "walk dir")
	checkf.temp = checkCmd.Flags().Bool("temp", false, "show contents of temporary storage struct for changes to user data")
	checkf.show = checkCmd.Flags().Bool("show", false, "show exe info")
	checkf.ask = checkCmd.Flags().Bool("ask", false, "check askconfirm")
}
