/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

const verstr string = "0.0.1"

func ce(e error) {
	if e != nil {
		panic(e)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotget",
	Short: "make groups of config files/dotfiles to copy to one central location",
	Long: `
Dotstrike is an application with the intent to make managing config files easier.
The current use is backing up config files from various locations into a single path/repo, and syncing in both directions.
It was primarly intended for Windows systems, where there is no designated common path 
or standard practice for storing these files. 

In practice, dotstrike is a simple file management tool that can group and sync files
and directories between the path where they are used and a storage/repo location.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if *version {
			cmd.Print("version: ", verstr)
		}
		if *pf.debug {
			cmd.Printf("DEBUG")
			gdump := dscore.GD.Dump()
			for _, l := range gdump {
				cmd.Println(l)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type persistentData struct {
	verbose, global, help *bool
	debug                 *bool
	src, tgt              *[]string // source and target currently set as Persistent Flags
}

// pf is the persistentData var that stores all persistent flag values
var pf persistentData
var version *bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	initfunctions := []func(){dscore.CoreConfig}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cobra.OnInitialize(initfunctions...) // pass all initialization functions here
	pf = persistentData{
		verbose: rootCmd.PersistentFlags().BoolP("verbose", "v", false, "shows additional details on execution (debug)"),
		global:  rootCmd.PersistentFlags().BoolP("global", "g", false, "target the global group"),
		src:     rootCmd.PersistentFlags().StringSliceP("source", "s", []string{}, "src"),
		tgt:     rootCmd.PersistentFlags().StringSliceP("target", "t", []string{}, "tgt"),
		debug:   rootCmd.PersistentFlags().Bool("debug", false, "debug"),
		// These seem to be defaults
		//help: rootCmd.PersistentFlags().BoolP("help", "?", false, "prints long help for command"),
	}
	// version is not default
	version = rootCmd.Flags().Bool("version", false, "print application version")
}
