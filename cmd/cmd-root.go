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

// TODO:(low-implement) SET UP LOGGING
func initLogging() {
	lf, e := pops.MakeOpenFileF("dslog.txt")
	if e != nil {
		panic(e)
	}
	logger := slog.NewTextHandler(lf, nil)
	_ = logger
}

type cmdWrapper struct {
	*cobra.Command
	args    []string
	specs   []*dscore.Spec
	runFunc func(*cobra.Command, []string)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotstrike",
	Short: "set up, configure, and run file copy jobs to/from specific paths or path groups",
	Long: `
Dotstrike is a file management tool that can group files/directories and sync them
	between source paths and one or more target destinations.

Super quick start:
To start, first create a spec:
	> ds spec "myspec"
From there, the spec needs a source (src) and a target (tgt)
	> ds src c:/my_files/
	> ds tgt 'd:/backups/personal files/'


Specs are the primary method of defining and storing copy job details. The spec command creates a new spec when provided with an alias.
The spec alias is the identifier for the spec; as such, it must be unique. All aliases are made lowercase, and stripped of spaces/tabs, forward slash/backslash, and at signs. Other symbols should be fine.

Specs contain sources (src) and targets (tgt). Sources are copied to Targets; in other words, a source points to a location containing files, and a target points to where you want those files to be copied to.
As your first spec has just been created, it will be automatically selected, and other operations will affect it directly.
If you have multiple specs, you will need to either:
	- select the spec you want to modify (using spec command with an existing spec name) before running other commands.
	- use the --spec flag at the end of the command to change selection for only that operation.
	
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		print("got this far")
		if *version {
			cmd.Print("version: ", verstr)
		}
		if *pFlags.debug {
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

// pFlags is the persistentData var that stores all persistent flag values
var pFlags persistentData
var version *bool

func (p *persistentData) componentFlags() {

}

//
//
/*
TODO:(HIGH-Refactor) REMOVE FLAGS THAT ARE NOT TRUE PERSISTENT, ADD TO COMMANDS WHERE USED.
	-> Actually, probably just remove every persistent flag (except verbose)
Reasons:
all: can't provide details of functionality from individual functions (I think?)
global: only maaybe used in spec or list or something
spec: doesn't work with current intent for run command
src/tgt: prob not useable for source/target commands
Flags to remove:
 - all
 - global
 - spec
 - debug

*/
func init() {
	cobra.OnInitialize(dscore.CoreConfig, dscore.InitTempData) // pass all initialization functions here
	cobra.OnFinalize(dscore.EndEncode)
	pFlags = persistentData{
		verbose: rootCmd.PersistentFlags().BoolP("verbose", "v", false, "shows additional details on execution"),
		all:     rootCmd.PersistentFlags().BoolP("all", "a", false, "applies command to 'all' applicable items (see command help for more detail)"),
		global:  rootCmd.PersistentFlags().BoolP("global", "g", false, "target the global group"), //uncertain, overlap with all?
		spec:    rootCmd.PersistentFlags().StringArrayP("spec", "s", nil, "spec"),
		src:     rootCmd.PersistentFlags().StringArrayP("source", "o", nil, "src"),
		tgt:     rootCmd.PersistentFlags().StringArrayP("target", "t", []string{}, "tgt"),
		// dev use
		debug: rootCmd.PersistentFlags().Bool("debug", false, "debug"),
	}
	pFlags.setup()
	// version is not default
	version = rootCmd.Flags().Bool("version", false, "print application version")
}
