/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

type runner struct {
	*cobra.Command
	specs                                        []*dscore.Spec
	args                                         []string
	flagY, flagNoSelectedSpec, flagAll           *bool
	flagNoFiles, flagAllDir                      *bool
	flagOverrides, flagRunPartial, flagManualRun *[]string
	bOverrides, bManual, bPartial                bool
	specNames                                    []string
}

var mainRun runner

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run spec copy job(s)",
	Long: `Run copy jobs for one or more specs.
Modify a run with one-time overrides, perform a partial run, or run a one time manually-entered run`,
	Run: mainRun.run,
}

func (r *runner) run(cmd *cobra.Command, args []string) {
	//Note: Using args for spec names? Don't need the --spec PersistentFlag then
	e := r.calculateBools()
	if e != nil {
		cmd.PrintErr(e)
		cmd.Print("\nearly terminate")
		return // no val func break
	}
	r.Command = cmd
	r.args = args

	if r.bManual {
		r.runManualJob()
		return
	}

	r.specs = r.makeSpecList()
	if len(r.specs) == 0 {
		cmd.Print("0 specs to run")
		return
	}

	conf := oneSpecOrUserConfirm("Run copy job for Specs ("+strings.Join(r.specNames, ", "), r.specs)
	if conf {
		r.prepAndRun()
	}
}

func (r *runner) makeSpecList() []*dscore.Spec {
	r.calculateBools()
	// reason for not writing directly to r.specs??
	specs := make([]*dscore.Spec, 0, len(r.args)+1)
	r.specNames = make([]string, 0, len(r.args))
	temp := dscore.TempData()
	if !*r.flagNoSelectedSpec {
		s := temp.SelectedSpec()
		specs = append(specs, s)
		r.specNames = append(r.specNames, s.Alias)
	}
	for _, alias := range r.args {
		s := temp.GetSpec(alias)
		if s != nil {
			specs = append(specs, s)
			r.specNames = append(r.specNames, s.Alias)
		} else {
			r.Printf("Spec '%s' not found\n", alias)
		}
	}
	return specs
}

func (r *runner) handleFlags() {

}

func (r *runner) prepAndRun() {

}

// calculates bManual, bPartial, bOverrides.
// Errors on unusable combination of flags/args
func (r *runner) calculateBools() error {
	estr := ""
	if manLen := len(*r.flagManualRun); manLen > 1 {
		r.bManual = true
	} else if manLen > 0 {
		estr = "Manual run flag requires 2 arguments minimum\n"
	}

	if prtLen := len(*r.flagRunPartial); prtLen > 0 && !r.bManual {
		r.bPartial = true
	} else if prtLen > 0 {
		estr += "Partial flag and Manual flag are mutually exclusive\n"
	}

	if len(*r.flagOverrides) > 0 {
		r.bOverrides = true
	}
	if estr == "" {
		return nil
	} else {
		return fmt.Errorf("%s", estr)
	}
}
func (r *runner) makeCopyJobs() {

}

func (r *runner) runManualJob() {

}

func init() {
	rootCmd.AddCommand(runCmd)
	mainRun.flagAll = runCmd.Flags().Bool("all-specs", false, "Run ALL spec copy jobs")
	mainRun.flagNoSelectedSpec = runCmd.Flags().Bool("no-selected", false, "Disable run of selected spec")
	mainRun.flagY = runCmd.Flags().BoolP("confirm", "y", false, "Auto-Confirm all prompts during run")
	mainRun.flagOverrides = runCmd.Flags().StringArray("override", []string{}, `Set one-time overrides with a space-separated list of 'prefName value' pairs; check spec help for more details on available options.`)
	mainRun.flagNoFiles = runCmd.Flags().BoolP("no-files", "n", false, "Disable filecopy for run. Use for dry runs, or in combination with --all-dir to copy only the directory structure")
	mainRun.flagAllDir = runCmd.Flags().BoolP("all-dirs", "d", false, "Copy all Source subdirectories, including empty subdirectories. Use with --no-files to only copy the directories themselves.")
	mainRun.flagRunPartial = runCmd.Flags().StringArray("partial", []string{}, "partial s=2 t=1")
	mainRun.flagManualRun = runCmd.Flags().StringArray("manual", []string{}, "manual ")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
