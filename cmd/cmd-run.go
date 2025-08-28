/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"maps"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"iidexic.dotstrike/dscore"
)

var allFlagNames []string

func init() {
	// Use Overrides? Or make flags for all?
	/*
		fOverridesUsage := `Set one-time overrides with a space-separated list of 'prefName value' pairs.
		check spec help for more details on available options.`
	*/
	mainRun.rtPrefs = make(map[dscore.ConfigOption]*bool)
	fNoFileUsage := "Disable filecopy for run. Use for dry runs, or with --all-dir to copy only the directory structure"
	fAllDirUsage := `Copy all Source subdirectories, including empty subdirectories. 
Use with --no-files to only copy the directories themselves.`
	fManualUsage := `--manual --src="srcpath,[additional paths]" --tgt="tgtpath,[additional paths]"
Use to run a one-time copy job. REQUIRES use of src and tgt flags to input paths to copy from/to. 
	Use override flag to set run configuration; by default current global prefs will be used. `
	fPartialUsage := `--partial="s(1,n),t(2,n)" 
Use to copy a selected subset of a given spec's sources and targets. To use, include indices, dir names, or aliases of sources and targets as shown`
	fGlobalTgtUsage := `Use "--GlobalTarget" to enable write to Global Target for all specs in run.
	Use --GlobalTarget="off" to forcibly disable write to GlobalTarget for full run, including specs that exclusively target the Global Target.`
	rootCmd.AddCommand(runCmd)
	mainRun.flagAll = runCmd.Flags().Bool("all-specs", false, "Run ALL spec copy jobs")
	mainRun.flagNoSelectedSpec = runCmd.Flags().Bool("no-selected", false, "Disable run of selected spec")
	mainRun.flagY = runCmd.Flags().BoolP("confirm", "y", false, "Auto-Confirm all prompts during run")
	//mainRun.flagOverrides = runCmd.Flags().StringArray("override", []string{}, fOverridesUsage)
	mainRun.rtPrefs[dscore.BoolNoFiles] = runCmd.Flags().BoolP("no-files", "n", false, fNoFileUsage)
	mainRun.rtPrefs[dscore.BoolCopyAllDirs] = runCmd.Flags().BoolP("all-dirs", "d", false, fAllDirUsage)
	mainRun.rtPrefs[dscore.BoolRootSubdir] = runCmd.Flags().Bool("make-subdir", false, "Makes a new dir in target folder to copy a spec into.\nDir is named with spec's alias if possible else numbers will be added")
	mainRun.rtPrefs[dscore.BoolSourceSubdirs] = runCmd.Flags().Bool("separate-sources", false, "Copies each source into a separate subdir; name is source's alias or source path's dir name.")

	mainRun.fOptGlobalTarget = runCmd.Flags().String("global-target", "", fGlobalTgtUsage)
	runCmd.Flag("global-target").NoOptDefVal = "on"

	mainRun.rtPrefs[dscore.BoolIgnoreRepo] = runCmd.Flags().Bool("ignore-repo", false, "add git repo to global ignores; Disables copying the .git dir")
	mainRun.rtPrefs[dscore.BoolIgnoreHidden] = runCmd.Flags().Bool("ignore-hidden", false, "add hidden paths to global ignores; Disables copy of paths that begin with `_` or `.`")
	mainRun.flagRunPartial = runCmd.Flags().StringArray("partial", []string{}, fPartialUsage)
	mainRun.fManualRun = runCmd.Flags().Bool("manual", false, fManualUsage)
	mainRun.flagSources = runCmd.Flags().StringArray("src", []string{}, `--src="path1,path2" (for partial/manual)`)
	mainRun.flagTargets = runCmd.Flags().StringArray("tgt", []string{}, `--tgt="path1,path2" (for partial/manual)`)

}

type runner struct {
	*cobra.Command
	specs []*dscore.Spec
	args  []string
	set   *pflag.FlagSet

	flagY, flagNoSelectedSpec, flagAll *bool
	rtPrefs                            map[dscore.ConfigOption]*bool
	realPrefs                          map[dscore.ConfigOption]bool
	fManualRun                         *bool
	fOptGlobalTarget                   *string

	flagOverrides, flagRunPartial *[]string
	flagSources, flagTargets      *[]string
	bOverrides, bPartial          bool
	manualMode, partialMode       bool
	specNames                     []string
	flagNames                     flagCatcher
}

type flagCatcher struct {
	used []string
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
	e := r.findFlags() //NOTE: Combine with handleFlags
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

func (r *runner) prepAndRun() {
	jm := dscore.JobManager()
	_ = jm

}

func (r *runner) processPartial() {

}

func (r *runner) makeSpecList() []*dscore.Spec {
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

// ── Flag/Config Logic ───────────────────────────────────────────────
func (r *runner) handleFlags() error {
	//nflags := r.set.NFlag() - len(r.rtPrefs)
	if len(r.rtPrefs) > 0 {
		r.processOptions()
	}

	if *r.fManualRun {

	} else if r.bPartial { // partial/manual should not be runing at same time
		r.processPartial()
	}

	return nil
}

// must be run for option use elsewhere.
// makes realPrefs (deref map)
func (r *runner) processOptions() {
	for k, v := range r.rtPrefs {
		r.realPrefs[k] = *v
	}

	// make globaltarget option from flag
	if optGT := *r.fOptGlobalTarget; optGT != "" {
		if toBool := dscore.StringToBool(optGT); toBool != nil {
			if *toBool {
				r.realPrefs[dscore.BoolUseGlobalTarget] = true
			} else if !*toBool { // this is 100% the only possibility but feel the need to be sure
				r.realPrefs[dscore.BoolKillGlobalTarget] = true
			}
		}
	}
}

func (r *runner) trimPrefs() {
	maps.DeleteFunc(r.realPrefs, r.keep)
}

func (r *runner) keep(k dscore.ConfigOption, v bool) bool { return k.IsRealOption() && k.IsBool() }

// ──────────────────────────────────────────────────────────────────────

// calculates bManual, bPartial, bOverrides.
// Errors on unusable combination of flags/args

func (r *runner) findFlags() error {
	estr := ""
	if !r.set.HasFlags() {
		return nil
	}
	r.set.Visit(r.flagNames.collect)

	if prtLen := len(*r.flagRunPartial); prtLen > 0 && !*r.fManualRun {
		r.bPartial = true // probably delete
	} else if prtLen > 0 { //if Partial flag + Manual flag?
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

func (fc *flagCatcher) collect(f *pflag.Flag) {
	fc.used = append(fc.used, f.Name)
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
