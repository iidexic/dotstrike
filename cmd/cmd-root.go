/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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

// // TODO:(low-implement) SET UP LOGGING
// func initLogging() {
// 	lf, e := pops.MakeOpenFileF("dslog.txt")
// 	if e != nil {
// 		panic(e)
// 	}
// 	logger := slog.NewTextHandler(lf, nil)
// 	_ = logger
// }

type cmdData struct {
	*cobra.Command
	args       []string
	specs      []*dscore.Spec
	components []*dscore.PathComponent
	ignoreptns []string
	countArgs  int
	msg        opString
	runFunc    func(*cobra.Command, []string)
}

type opString struct {
	operation, opVerb              string
	directType, parentType         string
	directNames, parentNames       string
	directAffected, parentAffected int
}

func (O opString) String() string {
	str := O.operation + "."
	if O.parentType != "" {
		str += fmt.Sprintf("This will %s %d %ss in %d %ss.", O.opVerb, O.directAffected, O.directType, O.parentAffected, O.parentType)
	} else {
		str += fmt.Sprintf("This will %s %d %ss.", O.opVerb, O.directAffected, O.directType)
	}
	if O.directNames != "" {
		str += fmt.Sprintf("\n%s %s Names: (%s)", O.opVerb, O.directType, O.directNames)
	}
	if O.parentNames != "" {
		str += fmt.Sprintf("\nAffected %s Names: (%s)", O.directType, O.directNames)
	}
	return str
}

func newCmdData(cmd *cobra.Command, args []string) *cmdData {
	c := &cmdData{
		args:      args,
		specs:     make([]*dscore.Spec, len(args)+1),
		countArgs: len(args),
	}
	c.Command = cmd
	return c
}

// Gets all specs from arglist. Returns selected spec if no arglist passed
// getSpecs gathers a list of indices from arglist that did NOT find a match
// It then does nothing with this and instead returns the qty of specs found?
func (C *cmdData) getSpecs(forceSelected bool, arglist ...string) int {
	var nf []int
	C.specs, nf = dscore.TempData().GetSpecs(forceSelected, arglist...)
	if nf == nil {
		return 1
	}
	return len(nf)
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
		if *version {
			cmd.Print("version: ", verstr)
		}
		if *persistentFlags.debug {
			cmd.Printf("DEBUG")
			gdump := dscore.DumpGlobals()
			for _, l := range gdump {
				cmd.Println(l)
			}
		}
		if !*persistentFlags.debug && !*version && len(args) == 0 {
			cmd.Print("add --help for usage details")
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
	verbose, all, debug *bool
	//global *bool
	//spec, src, tgt       *[]string // source and target currently set as Persistent Flags
	//bspec, bsrc, btgt    bool      //lazy
	countFlags int
}

func (p *persistentData) setup() {

	p.countFlags = 0 //just to make sure
	for _, b := range []bool{*p.verbose, *p.debug /*, *p.global, p.bspec, p.bsrc, p.btgt*/} {
		if b {
			p.countFlags++
		}
	}

}

// persistentFlags is the persistentData var that stores all persistent flag values
var persistentFlags persistentData
var version *bool

func configLoadInit() {
	e := dscore.CoreConfig()
	if e != nil {
		panic(e)
	}
}

func init() {
	cobra.OnInitialize(configLoadInit, dscore.InitTempData) // pass all initialization functions here
	cobra.OnFinalize(dscore.EndEncode)
	persistentFlags = persistentData{
		verbose: rootCmd.PersistentFlags().BoolP("verbose", "v", false, "shows additional details on execution"),
		all:     rootCmd.PersistentFlags().BoolP("all", "a", false, "applies command to 'all' applicable items (see command help for more detail)"),
		debug:   rootCmd.PersistentFlags().Bool("debug-secret", false, ""), //hide
	}
	rootCmd.PersistentFlags().MarkHidden("debug-secret")
	persistentFlags.setup()
	// version is not default
	version = rootCmd.Flags().Bool("version", false, "print application version")
}
