package dscore

import (
	"fmt"
	"maps"
	"strings"

	"iidexic.dotstrike/config"
	pops "iidexic.dotstrike/pathops"
)

//TODO: (hi-postrelease) Move job/execute code to pops package?
//	For the copy running, dscore is exclusively acting as middle man from cmd to pops
//	At the very least, find what makes more sense in pops and pull it over there.

var (
	ErrNotMade      = fmt.Errorf("Error: Spec Copy Jobs did not get made")
	ErrNoSpecConfig = fmt.Errorf("Error: Config not applied to spec")
)

var Copier = pops.Copier()

type jobProcessor struct { // uses gd/TempData() for global prefs and global target
	specs         map[string]*jobSpec
	runtimeConfig map[config.OptionKey]bool
	setupComplete bool
}

var manager = jobProcessor{
	specs:         make(map[string]*jobSpec),
	runtimeConfig: make(map[config.OptionKey]bool),
}

// Adds specs as partial specs by stripping any source/target
// paths that don't match with the IDs passed.
//
// The current implementation clones the specs, but will probably change
// to using specs directly from gd in the future.
func (J *jobProcessor) AddAsPartial(s *Spec, sourceIDs []string, targetIDs []string) error {
	// TODO: determine if we can just use specs direct from globalData as to not bloat things.
	isources := s.GetMatching(sourceIDs, true)
	itargets := s.GetMatching(targetIDs, false)
	ns := s.cloneSelf()
	ns.stripComponentList(isources, true)
	ns.stripComponentList(itargets, false)
	J.specs[ns.Alias] = &jobSpec{Spec: ns}
	return nil
}

func (J *jobProcessor) AddSpecs(s ...*Spec) {
	for i := range s {
		J.specs[s[i].Alias] = &jobSpec{Spec: s[i]}
	}
}

var setupRunContinue = false

func (J *jobProcessor) SetupAndRunAll(abortOnError bool) error {
	// if confirmHaveRuntimeConfig && (J.runtimeConfig == nil || len(J.runtimeConfig) == 0) {
	// 	return fmt.Errorf("Run terminated: No runtime config (confirmHaveRuntimeConfig on)")
	// }
	var runErr error
	for i, s := range J.specs {

		cerr := s.applyAndCheckConfigs(gd.data.Prefs.Bools, J.runtimeConfig)
		//HANDLE UseGlobalTarget, KillGlobalTarget before groups made
		if cerr != nil {
			return cerr
		}

		switch {
		case s.config[BoolUseGlobalTarget] && s.config[BoolKillGlobalTarget]:
			return fmt.Errorf("Confligting config options; UseGlobal, KillGlobal both true")
		case s.config[BoolUseGlobalTarget]:
			s.addGlobalTarget()
		case s.config[BoolKillGlobalTarget]:
			s.removeGlobalTarget()
		}

		s.group = Copier.NewJobGroup(J.specs[i].groupExport())
		J.setupComplete = true
		e := J.specs[i].group.RunAll(abortOnError)
		if e != nil {
			if abortOnError {
				return fmt.Errorf("Problem with run of group %s: %w", J.specs[i].Alias, e)
			}
			if runErr == nil {
				runErr = fmt.Errorf("Group Fails: (GRP-%s - %w)", J.specs[i].Alias, e)
			} else {
				runErr = fmt.Errorf("%w (GRP-%s - %w})", runErr, J.specs[i].Alias, e)
			}
		}

	}
	return runErr
}

func (J *jobProcessor) SetupOnly() error {
	// if confirmHaveRuntimeConfig && (J.runtimeConfig == nil || len(J.runtimeConfig) == 0) {
	// 	return fmt.Errorf("Run terminated: No runtime config (confirmHaveRuntimeConfig on)")
	// }
	var specErrors error
	errcount := 0
	for i := range J.specs {
		e := J.specs[i].applyAndCheckConfigs(gd.data.Prefs.Bools, J.runtimeConfig)
		if e != nil {
			errcount++
			if specErrors == nil {
				specErrors = fmt.Errorf("((  %w", e)
			} else {
				specErrors = fmt.Errorf("%w, %w", specErrors, e)
			}
		}
		J.specs[i].group = Copier.NewJobGroup(J.specs[i].groupExport())
	}

	if specErrors != nil {
		specErrors = fmt.Errorf("%d SPEC SETUP ERRORS: %w", errcount, specErrors)
		return specErrors
	}
	return nil
}

// TODO: (LOW-later) re-do configs so a priority can be attached (then only need to write overrides at spec-level)

// Performs necessary process of applying all priority levels of configuration to all individual specs
// Don't think this is needed anymore
func (J *jobProcessor) applySpecConfigs() {
	for i := range J.specs {
		//1. apply global config. Handles both defaults + globals
		maps.Copy(J.specs[J.specs[i].Alias].config, gd.data.Prefs.Bools)
		if J.specs[i].OverrideOn { //2. apply spec overrides where enabled
			maps.Copy(J.specs[J.specs[i].Alias].config, J.specs[i].Overrides.Bools)
		}
		if len(J.runtimeConfig) > 0 { // 3. apply runtimeConfig if populated
			maps.Copy(J.specs[J.specs[i].Alias].config, J.runtimeConfig)
		}
		J.specs[i].configApplied = true
	} // At this point, all specs have their configuration

}

func (J *jobProcessor) makeJobGroups() {
	for i := range J.specs {
		alias := J.specs[i].Alias
		J.specs[i].group = Copier.NewJobGroup(alias,
			J.specs[i].sourcePaths(),
			J.specs[i].targetPaths(),
			J.specs[i].config)
		J.specs[i].madeJobs = true
	}
}

func JobManager() *jobProcessor {
	//1. sort out config as is
	maps.Copy(manager.runtimeConfig, gd.data.Prefs.Bools)
	return &manager
}

// RuntimeConfigure Directly overwrites JobProcessor runtime config.
// This is the highest priority set of
// Remove any unnecessary key/value pairs before calling RuntimeConfigure.
func (J *jobProcessor) RuntimeConfigure(opts map[ConfigOption]bool) {
	maps.Copy(J.runtimeConfig, opts)
}

func (J *jobProcessor) SetupManual(sourcePaths, targetPaths []string) (*jobSpec, error) {
	s := Spec{Alias: "@manual@"}
	var e error
	ersrc := s.temporaryComponents(true, sourcePaths...)
	ertgt := s.temporaryComponents(false, targetPaths...)
	if ersrc != nil || ertgt != nil {
		e = fmt.Errorf("path add failures:\nsources:%w\ntargets:%w", ersrc, ertgt)
	}
	return &jobSpec{Spec: &s}, e
}

// func (J *jobProcessor) RunAll(stopOnError bool) { pops.Copier().RunAll(stopOnError) }

func (J *jobProcessor) WriteJobDetail() string {
	//len==
	dtl := make([]string, len(J.runtimeConfig)+len(J.specs)+2)
	i := 0
	dtl[i] = "COPIER\nConfig  for Copy Jobs:"
	i++
	// J.runtimeConfig should have Everything in it before now;
	for k, v := range J.runtimeConfig {
		dtl[i] = fmt.Sprintf("[%s] = %t", k.String(), v)
		i++
	}
	dtl[i] = "Job Specs:"
	i++
	for _, js := range J.specs {
		dtl[i] = js.briefDetail()
	}
	return strings.Join(dtl, "\n")
}
