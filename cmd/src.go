/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

type srcFlags struct {
	alias  *string
	ignore *[]string
}

var srcF = srcFlags{}

// srcCmd represents the src command
var srcCmd = &cobra.Command{
	Use:   "src",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		/*
			Alias   string        `toml:"alias"`
			Path    string        `toml:"path"`
			Ignores []string      `toml:"ignores"`
		*/
		// args = id for src (path, filepath.base(path), or alias)
		// what are affecting flags:
		if *pData.all {

		}
	},
}

func init() {
	//NOTE: does it make sense to have this attached to spec?
	// how tf

	//OK:
	// don't make a subcommand of spec; canot have a main cmd arg if running a subcmd
	// could make like poleless flags but sounds messy

	rootCmd.AddCommand(srcCmd)
	srcF.ignore = srcCmd.Flags().StringArray("ignore patterns", nil, "ignore")
	srcF.alias = srcCmd.Flags().String("Set Alias", "", "alias")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// srcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// srcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
