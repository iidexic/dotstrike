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

func init() {
	// Use Overrides? Or make flags for all?
	/*
		fOverridesUsage := `Set one-time overrides with a space-separated list of 'prefName value' pairs.
		check spec help for more details on available options.`
	*/
	fNoFileUsage := "Disable filecopy for run. Use for dry runs, or with --all-dir to copy only the directory structure"
	fAllDirUsage := `Copy all Source subdirectories, including empty subdirectories. 
Use with --no-files to only copy the directories themselves.`
	fManualUsage := `--manual --src="srcpath,[additional paths]" --tgt="tgtpath,[additional paths]"
Use to run a one-time copy job. REQUIRES use of src and tgt flags to input paths to copy from/to. 
	Use override flag to set run configuration; by default current global prefs will be used. `
	fPartialUsage := `--partial="s(1,n),t(2,n)"
provides ability to copy a subset of a given spec's sources and targets. To use, include indices or dir names of sources and targets as shown`
	rootCmd.AddCommand(runCmd)
	mainRun.flagAll = runCmd.Flags().Bool("all-specs", false, "Run ALL spec copy jobs")
	mainRun.flagNoSelectedSpec = runCmd.Flags().Bool("no-selected", false, "Disable run of selected spec")
	mainRun.flagY = runCmd.Flags().BoolP("confirm", "y", false, "Auto-Confirm all prompts during run")
	//mainRun.flagOverrides = runCmd.Flags().StringArray("override", []string{}, fOverridesUsage)
	mainRun.fOptNoFiles = runCmd.Flags().BoolP("no-files", "n", false, fNoFileUsage)
	mainRun.fOptAllDirs = runCmd.Flags().BoolP("all-dirs", "d", false, fAllDirUsage)
	mainRun.fOptGlobalTarget = runCmd.Flags().Bool("globaltarget", false, "'--globaltarget=false' to disable write to global target")
	mainRun.fOptNoRepo = runCmd.Flags().Bool("ignore-repo", false, "add git repo to global ignores; Disables copying .git dir/repo dir")
	mainRun.fOptNoHidden = runCmd.Flags().Bool("ignore-hidden", false, "add hidden paths  to global ignores; Disables copy of paths that begin with `_` or `.`")
	mainRun.flagRunPartial = runCmd.Flags().StringArray("partial", []string{}, fPartialUsage)
	mainRun.fOptManualRun = runCmd.Flags().Bool("manual", false, fManualUsage)
	mainRun.flagSources = runCmd.Flags().StringArray("src", []string{}, `--src="path1,path2" (for partial/manual)`)
	mainRun.flagTargets = runCmd.Flags().StringArray("tgt", []string{}, `--tgt="path1,path2" (for partial/manual)`)

}

type runner struct {
	*cobra.Command
	specs                              []*dscore.Spec
	args                               []string
	flagY, flagNoSelectedSpec, flagAll *bool

	fOptNoFiles, fOptAllDirs, fOptRootSubdir, fOptManualRun *bool
	fOptNoHidden, fOptNoRepo, fOptGlobalTarget              *bool

	flagOverrides, flagRunPartial *[]string
	flagSources, flagTargets      *[]string
	bOverrides, bPartial          bool
	specNames                     []string
}

var mainRun runner

// to handle different run modes/methods
type runMethod interface {
	execute() error
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run spec copy job(s)",
	Long: `Run copy jobs for one or more specs.
Modify a run with one-time overrides, perform a partial run, or run a one time manually-entered run`,
	Run: mainRun.run,
}

func (r *runner) run(cmd *cobra.Command, args []string) {
	// 1. check flags
	// 1a:
	// 1-2. make spec list
	e := r.calculateBools()
	if e != nil {
		cmd.PrintErr(e)
		cmd.Print("\nearly terminate")
		return // no val func break
	}
	r.Command = cmd
	r.args = args

	// if r.bManual {
	// 	r.runManualJob()
	// 	return
	// }

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

func (r *runner) handleFlags() error {
	// if r.bOverrides {
	//
	// 	var overrides map[string]bool = make(map[string]bool, 1)
	// }
	if *r.fOptManualRun {

	} else if r.bPartial { // partial/manual should not be runing at same time
		r.processPartial()
	}

	return nil
}

func (r *runner) prepAndRun() {
	jm := dscore.JobManager()
	_ = jm

}

func (r *runner) processPartial() {

}

// calculates bManual, bPartial, bOverrides.
// Errors on unusable combination of flags/args
func (r *runner) calculateBools() error {
	estr := ""

	if prtLen := len(*r.flagRunPartial); prtLen > 0 && !*r.fOptManualRun {
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

	if len(r.args) == 2 {
		jobmgr := dscore.JobManager()
		job, e := jobmgr.SetupManual([]string{r.args[0]}, []string{r.args[1]})
		_ = job
		if e != nil {
			r.PrintErr(e)
		} else {

		}

	}

}
