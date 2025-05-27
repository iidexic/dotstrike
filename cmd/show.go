/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show details of running program",
	Long:  `I'M LEARNING OVER HERE OKAY!!!`,
	Run: func(cmd *cobra.Command, args []string) {
		exec, e := os.Executable()
		ce(e)
		home, e := os.UserHomeDir()
		ce(e)
		cachedir, e := os.UserCacheDir()
		ce(e)
		cfgdir, e := os.UserConfigDir()
		ce(e)
		wd, e := os.Getwd()

		fmt.Print("-> show called\nExecutable:")
		fmt.Println(exec)
		fmt.Printf("~ = %s\n", home)
		fmt.Printf("cache dir: %s\n", cachedir)
		fmt.Println("config dir:", cfgdir)
		fmt.Printf("workingDir: %s\n", wd)
		fmt.Printf("dotstrike | VERSION", verstr)

	},
}

func checkflags(cmd *cobra.Command, flags []string) {

}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle"),``
}
