/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

type commandFlags struct {
	sourceFlag *[]string
	targetFlag *[]string
}

var flagValues = commandFlags{sourceFlag: nil, targetFlag: nil}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new config, floating source, or floating destination",
	Long: `used to make a new config component to store data on
	a set of source files/directories to be transferred to a target dir
	
	alternatively, add can be used with no alias to create floating (unassigned)
	source or target components for future use

	use (define config):
		> ds new [alias] <flags/args>
			alias - unique identifier used to reference this config in other commands.

	if alias is not provided, a source or target flag can still be passed
	components created in this way will be 'floating'(unassigned), and require an alias
	a source/target with an alias can be assigned to a component by alias instead of path
	example:
		>:ds new --target a='primary' ~/ds_storage
	(see flag descriptions for syntax details)

	flags:
		--source [path(s) or existing source alias] {a = 'alias'} 
			define one or more source paths, separated by a space.
		--target [path(s) or existing source alias] {a = 'alias'} 			
			define one or more target paths, separated by a space.
	`,
	Run: actionGroupAdd,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringArrayP("source", "s", []string{}, "source")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func actionGroupAdd(cmd *cobra.Command, args []string) {
	print("args: ")
	print(args)
	//check arg1 for valid alias
	if len(args) > 0 {

	}

}
