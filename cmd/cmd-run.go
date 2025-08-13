/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

type runner struct {
	*cobra.Command
	specs                              []*dscore.Spec
	args                               []string
	flagY, flagNoSelectedSpec, flagAll *bool
	flagNoFiles, flagAllDir            *bool
	flagOverrides                      *[]string
}

var mainRun runner

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the selected spec copy job",
	Long:  ``,
	Run:   mainRun.run,
}

func (r *runner) run(cmd *cobra.Command, args []string) {
	// 2 lines in and it's already wack, what even is this
	mainRun = runner{Command: cmd, flagY: r.flagY, flagOverrides: r.flagOverrides}
	r = &mainRun //lets just call this forced singleton pattern :)
	r.specs = r.makeSpecList()
}

func (r *runner) makeSpecList() []*dscore.Spec {
	listAlias := *pFlags.spec
	specs := make([]*dscore.Spec, 0, len(listAlias)+1)
	temp := dscore.TempData()
	if !*r.flagNoSelectedSpec {
		specs = append(specs, temp.SelectedSpec())
	}
	for _, alias := range listAlias {
		s := temp.GetSpec(alias)
		if s != nil {
			specs = append(specs, s)
		}
	}
	return specs
}

func init() {
	rootCmd.AddCommand(runCmd)
	mainRun.flagAll = runCmd.Flags().Bool("all", false, "Run ALL spec copy jobs")
	mainRun.flagNoSelectedSpec = runCmd.Flags().Bool("noselected", false, "Disable run of selected spec")
	mainRun.flagY = runCmd.Flags().BoolP("confirm", "y", false, "Auto-Confirm all prompts during run")
	mainRun.flagOverrides = runCmd.Flags().StringArray("override", []string{}, `Set one-time overrides with a space-separated list of 'prefName value' pairs; check spec help for more details on available options.`)
	mainRun.flagNoFiles = runCmd.Flags().BoolP("no-files", "n", false, "Disable filecopy for run. Use for dry runs, or in combination with --all-dir to copy only the directory structure")
	mainRun.flagAllDir = runCmd.Flags().BoolP("all-dirs", "d", false, "Copy all Source subdirectories, including empty subdirectories. Use with --no-files to only copy the directories themselves.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
