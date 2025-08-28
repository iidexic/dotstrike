package dscore

import (
	"fmt"
	"maps"

	"iidexic.dotstrike/config"
	pops "iidexic.dotstrike/pathops"
)

//TODO: Make Run command non-persistently modify spec prefs for hard overrides/one-time overrides

//todo list
// - finish hard overrides
// - implement manual runs
// - implement partial runsneovid
// - implement all flags

// Hard overrides;  what's  the joke (keeping this, was basically asleep when I wrote it)
// think thru hard override implementation:
/*

 */

type jobProcessor struct {
	specs         map[string]*jobSpec
	runtimeConfig map[config.OptionKey]bool
}

type jobSpec struct {
	*Spec
	config                  map[config.OptionKey]bool
	jobs                    []*pops.CopyJob
	manualRun               bool
	configApplied, madeJobs bool
}

var manager = jobProcessor{
	specs:         make(map[string]*jobSpec),
	runtimeConfig: make(map[config.OptionKey]bool),
}

func (J *jobProcessor) AddSpecs(s []*Spec) {
	//Preference priority hierarchy:
	/* lower prio to higher prio;
	0. program defaults
	1. global pref
	2. spec overrides.
	3. runtime flags/hard override
	In order for this to work:
	- Do not write global prefs to overrides. ENSURE THIS IS REMOVED
	- Do not write global prefs to runtime config. (remove, somewhere in this file)
	- There should be no reason to first write program defaults if init has occured correctly

	Order of Operations:
	1. directly apply global pref
	2. maps.Copy() spec overrides IF they are on
	3. maps.Copy() runtimeOverride if populated
	*/
	for _, spec := range s {
		J.specs[spec.Alias] = &jobSpec{Spec: spec}
		maps.Copy(J.specs[spec.Alias].config, gd.data.Prefs.Bools)
		if spec.OverrideOn {
			maps.Copy(J.specs[spec.Alias].config, spec.Overrides.Bools)
		}
		if len(J.runtimeConfig) > 0 {
			maps.Copy(J.specs[spec.Alias].config, J.runtimeConfig)
		}
	}
}

func JobManager() *jobProcessor {
	//1. sort out config as is
	manager.applyConfig(gd.data.Prefs.Bools)
	return &manager
}

// applyConfig directly overwrites job processor's runtime config, no checks are done beforehand
// Any configuring must be done before specs are loaded into jobProcessor
func (J *jobProcessor) applyConfig(bools map[config.OptionKey]bool) {
	maps.Copy(J.runtimeConfig, bools)
}

func (J *jobProcessor) SetupCopyJobs() {
	for k := range J.specs {
		J.specs[k].makeCopyJobs()
	}
}

// Configure applies the preferences provided in prefdata to the jobProcessor
// Any CopyJob run before processing ends will use these settings - they are prioritized over any other.
// returns a list of the prefdata keys that were NOT applied (do not match any config option)
func (J *jobProcessor) Configure(prefdata map[string]bool) []string {
	notFound := make([]string, 0, len(prefdata))
	for id, val := range prefdata {
		if opt := OptionID(id); opt != NotAnOption {
			J.runtimeConfig[opt] = val
		} else {
			notFound = append(notFound, id)
		}
	}
	return notFound
}

// NOTE: Directly overwrites JobProcessor runtime config;
// These will apply to everything in the run!
// Remove any unnecessary key/value pairs before calling RuntimeConfigure
func (J *jobProcessor) RuntimeConfigure(opts *map[ConfigOption]bool) {
	maps.Copy(J.runtimeConfig, *opts)

}

func (J *jobProcessor) SetupManual(sourcePaths, targetPaths []string) (*jobSpec, error) {
	s := Spec{Alias: "@manual@"}
	var e error
	ersrc := s.temporaryComponents(true, sourcePaths...)
	ertgt := s.temporaryComponents(true, sourcePaths...)
	if ersrc != nil || ertgt != nil {
		e = fmt.Errorf("path add failures:\nsources:%w\ntargets:%w", ersrc, ertgt)
	}
	return &jobSpec{Spec: &s}, e
}

func (js *jobSpec) editOverrides(runtimeConfig map[config.OptionKey]bool) {
	for k, v := range runtimeConfig {
		if k.IsBool() {
			js.OverrideOn = true
			js.Overrides.Bools[k] = v
		}
	}
}
func (js *jobSpec) makeCopyJobs() {
	for i, src := range js.Sources {
		for j, tgt := range js.Targets {
			nm := fmt.Sprintf("%s/%d/%d", js.Alias, i, j)
			js.jobs = append(js.jobs, pops.Copier().NewJob(nm, src.Path, tgt.Path))
		}
	}

}
func (js *jobSpec) RunJobs(continueOnError bool) error {
	if !js.madeJobs || !js.configApplied {
		return fmt.Errorf("jobs not initialized: (jobsMade=%t,configApplied=%t)",
			js.madeJobs, js.configApplied)
	}
	var ecopy []error
	for i := range js.jobs {
		e := js.jobs[i].Run()
		if e != nil && !continueOnError {
			return fmt.Errorf("err: copy from %s to %s\n%w", js.jobs[i].PathIn, js.jobs[i].PathOut, e)
		} else if e != nil {
			ecopy = append(ecopy, e)
		}
		//what to do here?
	}
	if len(ecopy) == 0 {
		return nil
	}
	e := fmt.Errorf("Copy errors: %d", len(ecopy))
	for i := range ecopy {
		e = fmt.Errorf("%w\n%w", e, ecopy[i])
	}
	return e
}

func (S *Spec) makeJobConfig() *prefs {
	return nil
}

// TODO: COMPLETELY RE-WRITE
func (S *Spec) applyJobConfig(job *pops.CopyJob) *pops.CopyJob {
	var runPrefs prefs
	if S.OverrideOn {
		runPrefs = S.Overrides
	} else {
		runPrefs = gd.data.Prefs
	}
	if runPrefs.Bools[BoolUseGlobalTarget] {
		//MakeSubdir
	}
	if runPrefs.Bools[BoolIgnoreRepo] {
		job.IgnoreGit()
	}

	return job
}

func (S *Spec) newCopyJob(isrc, itgt int) *pops.CopyJob {
	return pops.Copier().NewJob(
		S.jobName(isrc, itgt),
		S.Sources[isrc].Abspath,
		S.Targets[itgt].Abspath,
	)

}

func (S *Spec) jobName(isrc, itgt int) string {
	return fmt.Sprintf("%s.src-%d.tgt-%d", S.Alias, isrc, itgt)
}

type runtimeOverride struct {
	//*prefs //just add only changed to map
	options map[ConfigOption]bool
	on      bool
}

var hardOverrides = runtimeOverride{options: make(map[ConfigOption]bool), on: false} //unnecessary initialize?
// For runtime override flag use.
// var hardCopyOverride *prefs = &prefs{}
// var useHardOverride bool = false
//
// var overrideWhat []bool = make([]bool, len(BoolOptions))
//
// func MakeHardOverride() *prefs {
// 	// Add all options to the map, with their existing values in global (or spec if overrides already on)
// 	useHardOverride = true
// 	return hardCopyOverride
// }

// TODO:(NOW) - Re-write to use option.ValType IsBool
func SetHardOverride(boolOpt ConfigOption, value bool) bool {
	if boolOpt.IsBool() {
		hardOverrides.options[boolOpt] = value
		return true
	}

	return false
}
