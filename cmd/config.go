/* Copyright © 2025 NAME HERE <EMAIL ADDRESS> */
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	pops "iidexic.dotstrike/pathops"
)

// This will not be used
const evname string = "DOTSTRIKEROOT"
const cfgFile string = "dotstrikemainconfig.toml"

// Main config variable
var cfg = config{}

// config holds configuration status and data
type config struct {
	status configStatus
	cfpath string
	dpaths []string
	data   any
}

func _initCfg() {
	//2. loop through dpaths

	for _, p := range cfg.dpaths {
		fname := path.Join(p, cfgFile)
		print("[[Filepath:", fname, "]]")
		cf := pops.ReadF(fname)
		if !cf.Fail && len(cfg.cfpath) == 0 {

		}

	}
}

func defaultPaths(ddlist []string) {
	exec, e := os.Executable()
	ce(e)
	cachdir, e := os.UserCacheDir()
	ce(e)
	homedir, e := os.UserHomeDir()
	ce(e)
	ddlist = append(ddlist, exec)
	ddlist = append(ddlist, cachdir)
	ddlist = append(ddlist, homedir)
}

func coreConfig() {
	//check for config file, read, populate struct
	cfg.data = pops.MakeOpenFileF(cfg.cfpath)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
