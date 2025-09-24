/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"slices"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

var argstrSource = []string{"source", "sources", "src", "origin"}
var argstrTarget = []string{"target", "targets", "tgt", "destination"}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List user data",
	Long: `
	list all user data :)
`,
	Run: func(cmd *cobra.Command, args []string) {
		temp := dscore.TempData()
		switch {

		case len(args) == 0 && *persistentFlags.debug:
			cmd.Print("USER DATA ---------\n")
			cmd.Print("using datafile: ", dscore.ConfigTomlPath(), "\n")
			printsl := temp.Detail(*persistentFlags.verbose)
			for _, p := range printsl {
				cmd.Print(p, "\n")
			}
		case *persistentFlags.verbose:
			cmd.Print(temp.Detail(true))
		case len(args) > 0 &&
			(slices.Contains(argstrSource, args[0]) ||
				(slices.Contains(argstrTarget, args[0]))):
			cmd.Print(listComponents(dscore.TempData().Specs, true))
		case len(args) > 0:
			for _, a := range args {
				s := temp.GetSpec(a)
				if s != nil {
					cmd.Print(s.Detail())
				} else {
					cmd.Printf("spec '%s' not found", a)
				}
			}
		default:
			cmd.Println("User Specs:")
			specs := dscore.TempData().Specs
			for i := range specs {
				if dscore.TempData().Selected == i {
					cmd.Print("*** ", specs[i].ShortDetail(), " ***\n")
				} else {
					cmd.Println(specs[i].ShortDetail())
				}
			}
		}
	},
}

func listComponents(specs []dscore.Spec, isSource bool) []string {
	output := make([]string, 0, len(specs))
	for i := range specs {
		output = append(output, specs[i].DetailSources(true))
		output = append(output, specs[i].DetailTargets(true))
	}
	return output
}

func init() {
	rootCmd.AddCommand(listCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
