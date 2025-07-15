/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	pops "iidexic.dotstrike/pathops"
)

const verstr string = "0.0.1"

func ce(e error) {
	if e != nil {
		panic(e)
	}
}

// TODO: SET UP LOGGING
func initLogging() {
	lf, e := pops.MakeOpenFileF("dslog.txt")
	if e != nil {
		panic(e)
	}
	logger := slog.NewTextHandler(lf, nil)
	_ = logger
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotget",
	Short: "make groups of config files/dotfiles to copy to one central location",
	Long: `
Dotstrike is a tool to set up and trigger file copy actions.
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
		if *pData.debug {
			cmd.Printf("DEBUG")
			gdump := dscore.DumpGlobals()
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
	verbose, global, all *bool
	debug                *bool
	spec, src, tgt       *[]string // source and target currently set as Persistent Flags
	bspec, bsrc, btgt    bool      //lazy
	countFlags           int
}

func (p *persistentData) setup() {
	p.bsrc = len(*p.src) > 0
	p.bspec = len(*p.spec) > 0
	p.btgt = len(*p.tgt) > 0
	p.countFlags = 0 //just to make sure
	for _, b := range []bool{*p.verbose, *p.debug, *p.global, p.bspec, p.bsrc, p.btgt} {
		if b {
			p.countFlags++
		}
	}

}

type pfid int

//func (p *persistentData) checkAddData() { }

// pData is the persistentData var that stores all persistent flag values
var pData persistentData
var version *bool

func (p *persistentData) componentFlags() {

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cobra.OnInitialize(dscore.CoreConfig, dscore.InitTempData, dscore) // pass all initialization functions here
	cobra.OnFinalize(dscore.EndEncode)
	pData = persistentData{
		// TODO: determine whether verbose is a cobra built-in flag, or if there are other builtin besides help
		verbose: rootCmd.PersistentFlags().BoolP("verbose", "v", false, "shows additional details on execution (debug)"),
		all:     rootCmd.PersistentFlags().BoolP("all", "a", false, "applies command to 'all' applicable items (see command help for more detail)"),
		global:  rootCmd.PersistentFlags().BoolP("global", "g", false, "target the global group"), //uncertain, overlap with all?
		//NOTE: StringArrayP REQUIRES at least one flag argument
		// make sure this is acceptable for all use cases.
		// I would prefer them to also function as bools
		spec: rootCmd.PersistentFlags().StringArrayP("cfg", "c", nil, "cfg"),
		src:  rootCmd.PersistentFlags().StringArrayP("source", "s", nil, "src"),
		tgt:  rootCmd.PersistentFlags().StringArrayP("target", "t", []string{}, "tgt"),
		// dev use
		debug: rootCmd.PersistentFlags().Bool("debug", false, "debug"),
		// Help is default/built-in
		//help: rootCmd.PersistentFlags().BoolP("help", "?", false, "prints long help for command"),
	}
	pData.setup()
	// version is not default
	version = rootCmd.Flags().Bool("version", false, "print application version")
}
