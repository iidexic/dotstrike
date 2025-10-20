/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"iidexic.dotstrike/dscore"
)

var ( //TODO: (low) simplify with map/etc.
	// ── config.OptionKeys: ──────────────────────────
	NotAnOption               = dscore.NotAnOption
	bNoRepo, bNoHidden        = dscore.BoolIgnoreRepo, dscore.BoolIgnoreHidden
	bRootSubdir               = dscore.BoolRootSubdir
	bNoFiles, bAllDirs        = dscore.BoolNoFiles, dscore.BoolCopyAllDirs
	bUseGlobTgt, bKillGlobTgt = dscore.BoolUseGlobalTarget, dscore.BoolKillGlobalTarget
	configIDs                 = []dscore.ConfigOption{bNoRepo, bNoHidden, bRootSubdir,
		bNoFiles, bAllDirs, bUseGlobTgt, bKillGlobTgt} // bSeparate

	// ── Non-config pkg flags: ───────────────────────
	flagnameAll      = "all-specs"
	flagnameSelected = "selected"
	flagnameY        = "confirm"
	flagnamePartial  = "partial"
	flagnameManual   = "manual"
	flagnameSrc      = "src"
	flagnameTgt      = "tgt"
	//flagnameGlobTgt    = "global-target"
	flagnameNoRunDebug = "setup-only-debug"
	flagnameQuiet      = "quiet"
)

var (
	ErrMultiMode          = fmt.Errorf("cannot run multiple modes at once (standard, manual, or partial)")
	ErrMissModeDependency = fmt.Errorf("Missing Required flags/args for the mode selected!")
)

func flagOptKey(flagName string) dscore.ConfigOption {
	for _, opt := range configIDs {
		if flagName == opt.NameFlag() {
			return opt
		}
	}
	return dscore.NotAnOption
}

func init() {
	mainRun.rtPrefs = make(map[dscore.ConfigOption]*bool)
	mainRun.set = runCmd.Flags()
	rootCmd.AddCommand(runCmd)

	runMakeFlags()
}

// adds all runCmd flags to the command's flagset, including config flags
func runMakeFlags() {
	fManualUsage := `--manual --src="srcpath,[additional paths]" --tgt="tgtpath,[additional paths]"
Use to run a one-time copy job. REQUIRES use of src and tgt flags to input paths to copy from/to.
Use override flag to set run configuration; by default current global prefs will be used. `

	fPartialUsage := `--partial --src="srcpath/basedir/id,[additional paths]" --tgt="tgtpath,[additional paths]"
Use to run only specified sources/targets from one spec. REQUIRES use of src and tgt flags to specify copy job paths.
Use the override flag to set run configuration; by default current global prefs will be used.`

	f := runCmd.Flags()

	mainRun.flagAll = f.Bool(flagnameAll, false, "Run ALL spec copy jobs")
	mainRun.flagSelected = f.Bool(flagnameSelected, false, "Add selected spec to the run (if not already included)")
	mainRun.flagY = f.BoolP(flagnameY, "y", false, "Auto-Confirm all prompts during run")
	mainRun.fPartialRun = f.Bool(flagnamePartial, false, fPartialUsage)
	mainRun.fManualRun = f.Bool(flagnameManual, false, fManualUsage)
	mainRun.flagSources = f.StringArray(flagnameSrc, []string{}, `--src="path1,path2" for manual run;  --src="0,1,alias" for partial run (source index, alias, or dirname)`)

	mainRun.flagTargets = f.StringArray(flagnameTgt, []string{}, `--tgt="path1,path2" for manual run;  --tgt="0,1,alias" for partial run (target index, alias, or dirname)`)
	mainRun.fSetupDebug = f.Bool(flagnameNoRunDebug, false, "")
	mainRun.fQuiet = f.BoolP(flagnameQuiet, "q", false, "-q or --quiet suppresses all output")
	mainRun.set.MarkHidden(flagnameNoRunDebug)
	// mainRun.fAllToGlobalTarget = runCmd.Flags().String(flagnameGlobTgt, "", bUseGlobTgt.RunUsage())
	// runCmd.Flag("global-target").NoOptDefVal = "on"
	initConfigFlags()
}

// creates all flags in configIDs. Data is stored in config package for now.
func initConfigFlags() {
	for _, opt := range configIDs {
		if ns := opt.NameFshort(); ns != "" && len(ns) == 1 {
			mainRun.rtPrefs[opt] = mainRun.set.BoolP(opt.NameFlag(), ns, false, opt.RunUsage())
			continue
		}
		mainRun.rtPrefs[opt] = mainRun.set.Bool(opt.NameFlag(), false, opt.RunUsage())
	}
}

type runner struct {
	*cobra.Command
	specs []*dscore.Spec
	args  []string
	set   *pflag.FlagSet

	rtPrefs                      map[dscore.ConfigOption]*bool
	FinalConfig                  map[dscore.ConfigOption]bool
	flagY, flagSelected, flagAll *bool // checked where used
	fManualRun, fPartialRun      *bool // check first to toggle operation
	fAllToGlobalTarget           *bool // must convert into FinalConfig if used
	fSetupDebug, fQuiet          *bool
	runTriggered, dbg, setupOnly bool

	//flagOverrides, flagRunPartial *[]string
	flagSources, flagTargets *[]string
	manualMode, partialMode  bool
	specNames                []string
	flagsPassed              []string
}

var mainRun runner

func detailMainRun() string {
	detail := make([]string, 0, 32)
	hpref := "----[  PREFS: ]----\n-- RT bools --"
	rtp := ""
	for k, v := range mainRun.rtPrefs {
		rtp += fmt.Sprintf("(%s) = %t\n", k.String(), *v)
	}
	sep := "--------"
	hoflag := "-- Other Flags --"
	detail = append(detail, hpref, sep, rtp, hoflag)
	for i, f := range mainRun.flagsPassed {
		detail = append(detail, fmt.Sprintf("[%d] flag: %s", i, f))
	}
	// oo := ""
	//
	// runCmd.Flags().Visit(func(f *pflag.Flag) {
	// 	oo = "VISITED A FLAG"
	// 	oo += fmt.Sprintf("(%s: %v)", f.Name, f.Value)
	//
	// 	if len(oo)%3 == 0 {
	// 		oo += "\n"
	// 	}
	// })
	// detail = append(detail, hoflag, oo, sep, "--Specs--")
	// detail = append(detail, mainRun.specNames...)
	//
	return strings.Join(detail, "\n")
}

// to handle different run modes/methods
type runMethod interface {
	execute() error
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run spec copy job(s)",
	Long: `Run copy jobs for one or more specs.
Modify a run with one-time overrides, perform a partial run, or run a one-time, manually-entered run`,
	RunE: mainRun.run,
}

// this function gotta get cleaned up

func (r *runner) run(cmd *cobra.Command, args []string) error {
	//Intended to prevent risk of destroying user data with an accidental encode/write to file
	defer func() { dscore.TempData().Modified = false }()
	r.dbg = *r.fSetupDebug
	r.Command = cmd
	r.args = args
	JM := dscore.JobManager()
	_ = JM
	e := r.makeRuntimeConfig() //always returns nil error for now.

	if e != nil {
		return e
	}

	if r.dbg {
		cmd.Println("Running run's run function")
		cmd.Println("Runner pre-flag:")
		cmd.Printf("%s", r.dbgOut())
	}
	r.handleFlags()
	e = r.checkFlagProblems()
	if e != nil {
		return e
	}
	if r.dbg {
		cmd.Printf("Runner post-flag:\n%s", r.dbgOut())
	}

	if r.manualMode {
		e = r.runManualJob()
		// Need to do anything at the end?
		return e
	} else {
		r.specs = r.makeSpecList()
	}

	switch {
	case len(r.specs) == 0:
		cmd.Println("0 specs to run")
		return fmt.Errorf("No specs to run")
	case r.partialMode:
		if r.dbg {
			cmd.Println("Running Partial")
		}
		e := r.processPartial()
		if e != nil {
			return e
		}
	default:
		if r.dbg {
			cmd.Println("adding specs to JobManager")
		}
		dscore.JobManager().AddSpecs(r.specs...) // load specs
	}

	//jobDetail:= dscore.JobManager.MakeDetail()
	//conf := oneSpecOrUserConfirm("Run copy job for Specs ("+strings.Join(r.specNames, ", "), r.specs)
	conf := r.jobDetailUserConfirm()
	if conf {
		if r.dbg {
			cmd.Println("prepping and running")
		}
		err := r.prepAndRun()
		if err != nil {
			return fmt.Errorf("in prepAndRun(): %w", err)
		}
	}

	mainRun.verboseRunOutput()
	return nil
}

func (r *runner) jobDetailUserConfirm() bool {
	ls := len(r.specs)
	requestText := dscore.JobManager().WriteJobDetail()
	return ls == 1 || (ls > 1 && checkConfirm(requestText, r.flagY))

}

func (r *runner) prepAndRun() error {
	jm := dscore.JobManager()
	jm.RuntimeConfigure(r.FinalConfig)
	jm.AddSpecs(r.specs...)
	if r.dbg {
		r.Printf("JobManager: added specs, set runtimeConfig:\n%v\n", r.FinalConfig)
		r.Printf("JobManager Pre-Setup State:\n%+v\n", *jm)
		r.Println("running setup only")
		e := jm.SetupOnly()
		if e != nil {
			return e
		}
		r.Printf("JobManager Setup State:\n%+v\n", *jm)
		return nil
	} else {
		r.runTriggered = true
		// TODO: Now that there are Spec Setup Errors, split those into a different function.
		//	Then will be able to react to  spec failures
		err := jm.SetupAndRunAll(true)
		return err
	}
}

func (r *runner) finishRun() error {
	r.runTriggered = true
	err := dscore.JobManager().SetupAndRunAll(true)
	return err
}

func (r *runner) processPartial() error {
	J := dscore.JobManager()
	J.RuntimeConfigure(r.FinalConfig)
	for _, s := range r.specs {
		// keepSources := s.GetMatching(*r.flagSources, true) //keepTargets := s.GetMatching(*r.flagTargets, false)
		e := J.AddAsPartial(s, *r.flagSources, *r.flagTargets)
		if e != nil {
			return e
		}
	}
	return nil
}

func (r *runner) makeSpecList() []*dscore.Spec {
	// reason for not writing directly to r.specs??
	specs := make([]*dscore.Spec, 0, len(r.args)+1)
	r.specNames = make([]string, 0, len(r.args)+1)
	temp := dscore.TempData()
	for _, alias := range r.args {
		s := temp.GetSpec(alias)
		if s != nil {
			specs = append(specs, s)
			r.specNames = append(r.specNames, s.Alias)
		} else {
			r.Printf("Spec '%s' not found\n", alias)
		}
	}

	if ss := temp.SelectedSpec(); (*r.flagSelected || len(r.args) == 0) && !slices.Contains(specs, ss) {
		specs = append(specs, ss)
		r.specNames = append(r.specNames, ss.Alias)
	}
	if r.dbg {
		r.Printf("Specs pulled:\n,")
		for i := range specs {
			fmt.Printf("	[%d] %s", i, specs[i].Alias)
		}
	}
	return specs
}
func (r runner) verboseRunOutput() {
	if *persistentFlags.verbose {
		r.Print(dscore.Copier.GroupDetails())
	}
}

// ── Flag/Config Logic ───────────────────────────────────────────────

func (r *runner) makeRuntimeConfig() error {
	if !r.set.HasFlags() {
		return nil
	}
	r.set.Visit(r.process)
	return nil
}

// Processes flags and finds names.
// will make realPrefs only keeping input flags
func (r *runner) process(f *pflag.Flag) {
	r.flagsPassed = append(r.flagsPassed, f.Name)
	if key := flagOptKey(f.Name); key.IsBool() {
		r.FinalConfig[key] = *r.rtPrefs[key]
	}
}

//TODO:(mid) Take a step back; determine if there is a way besides loop Visit({get names}) then loop(names){if = loop(flagOptIDs) then add 2 finalconfig}

// handleFlags applies the remainder of flags that aren't checked in-process.
// Flags not handled by handleFlags+makeRuntimeConfig: All, NoSelect, Src, Tgt, Y
func (r *runner) handleFlags() {
	// if optGT := *r.fAllToGlobalTarget; optGT != "" { if toBool := dscore.StringToBool(optGT); toBool != nil { if *toBool {
	// 			r.FinalConfig[dscore.BoolUseGlobalTarget] = true } else { r.FinalConfig[dscore.BoolKillGlobalTarget] = true } } }
	r.manualMode = *r.fManualRun
	r.partialMode = *r.fPartialRun
}

func (r *runner) checkFlagProblems() error {
	if r.manualMode && r.partialMode {
		return ErrMultiMode
	}
	if r.manualMode || r.partialMode {
		if len(r.args) != 2 && (len(*r.flagSources) == 0 || len(*r.flagTargets) == 0) {
			return ErrMissModeDependency
		}
	}
	if *r.rtPrefs[bUseGlobTgt] && *r.rtPrefs[bKillGlobTgt] {

	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────

func (r *runner) runManualJob() error {

	// do we actually want operation using args in manualMode
	if len(r.args) == 2 {
		jobmgr := dscore.JobManager()
		job, e := jobmgr.SetupManual([]string{r.args[0]}, []string{r.args[1]})
		_ = job
		if e != nil {
			r.PrintErr(e)
		} else {
			jobmgr.RuntimeConfigure(r.FinalConfig)
		}

	}

	return nil
}

// *cobra.Command
// specs []*dscore.Spec
// args  []string
// set   *pflag.FlagSet
//
// rtPrefs                      map[dscore.ConfigOption]*bool
// FinalConfig                  map[dscore.ConfigOption]bool
// flagY, flagSelected, flagAll *bool   // checked where used
// fManualRun, fPartialRun      *bool   // check first to toggle operation
// fOptGlobalTarget             *string // must convert into FinalConfig if used
// fSetupDebug                  *bool
// dbg                          bool
//
// //flagOverrides, flagRunPartial *[]string
// flagSources, flagTargets *[]string
// manualMode, partialMode  bool
// specNames                []string
// flagsPassed              []string

func (r *runner) dbgOut() string {
	d := "[Runner]\n"
	d += fmt.Sprintf("args:'%s'\n", r.args)
	if len(r.specs) > 0 {
		d += "specs:\n"
		for i := range r.specs {
			d += fmt.Sprintf("	%s\n", r.specs[i].Alias)
		}
	}
	if len(r.flagsPassed) > 0 {
		d += "flags passed:\n"
		for _, f := range r.flagsPassed {
			d += fmt.Sprintf("	%s\n", f)
		}
	}
	return d
}
