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
	show, temp, walk, ask, path *bool
	parray, pslice              *[]string
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
		switch {
		case *checkf.ask:
			conf := askConfirmf("Is this question true")
			fmt.Printf("function gave: %t\n", conf)
			confright := askConfirmf("Was that correct")
			if confright {
				print("thats good")
			} else {
				print("oh no")
			}
		case len(*checkf.pslice) > 0:
			cmd.Printf("len pslice = %d\n", len(*checkf.pslice))
			cmd.Printf("len parray = %d\n", len(*checkf.parray))
			cmd.Printf("len args = %d\n", len(args))
			cmd.Println("pslice args:\n-----------")
			for i, str := range *checkf.pslice {
				cmd.Printf("[%d] %s\n", i, str)
			}
			printArgs(cmd, args)
		case len(*checkf.parray) > 0:
			cmd.Printf("len pslice = %d\n", len(*checkf.pslice))
			cmd.Printf("len parray = %d\n", len(*checkf.parray))
			cmd.Println("parray args:\n-----------")
			for i, str := range *checkf.parray {
				cmd.Printf("[%d] %s\n", i, str)
			}
			cmd.Println("-----------")
			printArgs(cmd, args)
		default:
			checkf.oldDefaultBehavior(cmd, args)
		}

	},
}
var checkf = checkCmdFlags{}

func printArgs(cmd *cobra.Command, args []string) {
	cmd.Printf("len args = %d\nargs:\n-----------\n", len(args))
	for i, a := range args {
		cmd.Printf("[%d] %s\n", i, a)
	}
}

func (c *checkCmdFlags) oldDefaultBehavior(cmd *cobra.Command, args []string) {
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

func init() {
	rootCmd.AddCommand(checkCmd)
	checkf.walk = checkCmd.Flags().BoolP("walk", "w", false, "walk dir")
	checkf.temp = checkCmd.Flags().Bool("temp", false, "show contents of temporary storage struct for changes to user data")
	checkf.path = checkCmd.Flags().BoolP("path", "p", false, "debug args path processing")
	// StringArray would be set multiple times; one arg per flag
	checkf.parray = checkCmd.Flags().StringArray("parray", []string{}, "print farg array")
	checkf.pslice = checkCmd.Flags().StringSlice("pslice", []string{}, "print farg slice")
	checkf.show = checkCmd.Flags().Bool("show", false, "show exe info")
	checkf.ask = checkCmd.Flags().Bool("ask", false, "check askconfirm")
}
