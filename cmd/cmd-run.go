/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

type runner struct {
	*cobra.Command
	specs                     []*dscore.Spec
	args                      []string
	flagY, flagNoSelectedSpec *bool
	flagOverrides             *[]string
}

var mainRun runner

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the selected spec copy job",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func (r *runner) runRun(cmd *cobra.Command, args []string) {
	// 2 lines in and it's already wack, what even is this
	mainRun = runner{Command: cmd, flagY: r.flagY, flagOverrides: r.flagOverrides}
	r = &mainRun //lets just call this forced singleton
	nSpecsEst := 1
	if pFlags.bspec {
		nSpecsEst += len(*pFlags.spec)
	}
	r.specs = r.makeSpecList(nSpecsEst)
}

func (r *runner) makeSpecList(n int) []*dscore.Spec {
	listAlias := *pFlags.spec
	specs := make([]*dscore.Spec, 0, n)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
