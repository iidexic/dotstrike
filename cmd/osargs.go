/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// osargsCmd represents the osargs command
var osargsCmd = &cobra.Command{
	Use:   "osargs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// os.Args will be:
		//0 = exe path i.e. 'c:\dev\bin\ds.exe'
		//	NOT the caller exe path or cwd
		// 1+ = the commands or arg passed past exe-name
		for i, s := range os.Args {
			println(i, ":", s)
		}
	},
}

func init() {
	showCmd.AddCommand(osargsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// osargsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// osargsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
