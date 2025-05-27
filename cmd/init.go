/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

type initLocator int

const (
	execPath initLocator = iota
	cacheDirADLocal
	homeDir
)

type initializeSettings struct {
	dir      *string
	priority initLocator
}

var iset initializeSettings

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "set up dotstrike's main config file and method of locating that config file",
	Long: `set up dotstrike's main required configuration.
	This includes both a main config file, as well as a file or environment variable that points to that file's directory.
	standard filepath: '~/.config/dotstrike/dsglobal.toml'
	possible path identifiers:
		2. 'dotstrikeroot.toml' file in same directory as dotstrike executable
		3. 'dotstrikeroot.toml' file in {cachedir}/dotstrike dir (Appdata/Local/dotstrike on Windows)
	
	If none of these are detected, running init attempts to create them, in the priority listed.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		print("checking for config")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	iset.dir = initCmd.Flags().String("directory", "~/.config/dotstrike/dsglobal.toml", "dir")
	pref := initCmd.Flags().String("preferred path identifier", "env", "prefer")
	switch *pref {
	case "exec", "exe", "executable directory", "exe dir", "2":
		iset.priority = execPath
	case "cache", "cache dir", "cache directory", "appdata", "appdata/local", "3":
		iset.priority = cacheDirADLocal
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
