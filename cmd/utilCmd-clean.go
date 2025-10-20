/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Use to clean up malformed specs or specs with duplicate aliases",
	Long: `Use to clean up improperly aliased specs.

All specs with empty, blank, or space-only aliases will be marked for deletion.
If multiple specs have the same alias, it is indeterminate which will be deleted and which will be kept.
If "test" or "temp" is passed as the only argument, aliases with that prefix will also be marked for deletion.

User will then be prompted to confirm deletion of each alias.
`,
	Run: func(cmd *cobra.Command, args []string) {
		mapSpecs := make(map[string]*dscore.Spec)
		specs := dscore.TempData().Specs
		deleteList := make([]string, 0, len(specs))
		var a string
		if len(args) > 0 {
			a = args[0]
		}
		istt := a == "test" || a == "temp"

		for i := range specs {
			_, ok := mapSpecs[specs[i].Alias]
			if ok {
				deleteList = append(deleteList, specs[i].Alias)
				continue
			}
			if strings.TrimSpace(specs[i].Alias) == "" {
				deleteList = append(deleteList, specs[i].Alias)
				continue
			}
			if istt {
				if strings.HasPrefix(specs[i].Alias, a) {
					deleteList = append(deleteList, specs[i].Alias)
					continue
				}
			}

			mapSpecs[specs[i].Alias] = &specs[i]
		}
		cmd.Printf("%d specs marked for deletion.\n Specs:", len(deleteList))
		for _, s := range deleteList {
			cmd.Printf(" '%s',", s)
		}
		cmd.Println("These specs either have no alias or a duplicate alias")
		skip := cleanDelete(cmd, deleteList)
		if len(skip) > 0 {
			cmd.Printf(`Skipped %d specs.
User data file is still malformed and error-prone. Please manually edit.`, len(skip))
		} else {

		}
	},
}

func cleanDelete(cmd *cobra.Command, deleteList []string) []string {
	skip := make([]string, 0, len(deleteList))
	y := *confirm
	td := dscore.TempData()
	for _, s := range deleteList {
		spec := td.GetSpec(s)
		if y || askConfirmf("Delete spec '%s' (%d Sources, %d Targets)", s, len(spec.Sources), len(spec.Targets)) {

			del := dscore.TempData().DeleteSpec(s)
			if del {
				cmd.Println("Deleted spec")
			} else {
				cmd.Println("Failed to delete spec")
			}
		} else {
			skip = append(skip, s)
			cmd.Println("Skipped spec")
		}
	}
	return skip
}

func cleanRename(cmd *cobra.Command, skipList []string) {

}

var confirm *bool

func init() {
	rootCmd.AddCommand(cleanCmd)
	f := cleanCmd.Flags()
	confirm = f.BoolP("confirm", "C", false, "confirm deletion")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
