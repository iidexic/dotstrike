package dscore

import (
	"fmt"
	"maps"

	"iidexic.dotstrike/config"
	pops "iidexic.dotstrike/pathops"
	"iidexic.dotstrike/uout"
)

//TODO: (hi-postrelease) Move job/execute code to pops package?
//	At the very least, find what makes more sense in pops and pull it over there.

var (
	ErrNotMade      = fmt.Errorf("Error: Spec Copy Jobs did not get made")
	ErrNoSpecConfig = fmt.Errorf("Error: Config not applied to spec")
)

var Copier = pops.Copier()

type jobProcessor struct { // uses gd/TempData() for global prefs and global target
	specs         map[string]*jobSpec
	runtimeConfig config.ConfigMap
	setupComplete bool
}

var manager = jobProcessor{
	specs:         make(map[string]*jobSpec),
	runtimeConfig: make(config.ConfigMap),
}

// Adds specs as partial specs by stripping any source/target
// paths that don't match with the IDs passed.
//
// The current implementation clones the specs.
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

func (J *jobProcessor) SetupAndRunAll(abortOnError bool) error {
	// if confirmHaveRuntimeConfig && (J.runtimeConfig == nil || len(J.runtimeConfig) == 0) {
	// 	return fmt.Errorf("Run terminated: No runtime config (confirmHaveRuntimeConfig on)")
	// }
	var runErr error
	for i, s := range J.specs {

		cerr := s.applyAndCheckConfigs(gd.data.Prefs.Bools, J.runtimeConfig)
		if cerr != nil {
			return cerr
		}

		// HANDLE UseGlobalTarget, KillGlobalTarget before groups made
		switch {
		case s.config[BoolUseGlobalTarget] && s.config[BoolKillGlobalTarget]:
			return fmt.Errorf("Confliggting config options; UseGlobal, KillGlobal both true")
		case s.config[BoolUseGlobalTarget]:
			s.addGlobalTarget()
		case s.config[BoolKillGlobalTarget]:
			s.removeGlobalTarget()
		}

		s.group = Copier.NewJobGroup(J.specs[i].groupExport())
		s.group.ConfigToJobs()
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
		J.specs[i].group.ConfigToJobs()
	}

	if specErrors != nil {
		specErrors = fmt.Errorf("%d SPEC SETUP ERRORS: %w", errcount, specErrors)
		return specErrors
	}
	return nil
}

func JobManager() *jobProcessor {
	//bug: Copying to runtimeConfig breaks priority. Removed Copy()
	if manager.runtimeConfig == nil {
		manager.runtimeConfig = make(map[ConfigOption]bool)
	}
	return &manager
}

func (J jobProcessor) String() string {
	out := uout.NewOut("==[JOB MANAGER]==")
	out.F("Setup Complete: %t", J.setupComplete)
	out.V("Job Specs:")
	out.IndR()
	for name, js := range J.specs {
		out.F("%s: %s", name, js.String())
	}

	return out.String()
}

// RuntimeConfigure Directly overwrites JobProcessor runtime config.
// This is the highest priority set of prefs?
// ---> Remove any unnecessary key/value pairs before calling RuntimeConfigure. (what)
func (J *jobProcessor) RuntimeConfigure(opts map[ConfigOption]bool) {
	if len(opts) > 0 {
		maps.Copy(J.runtimeConfig, opts)
	}
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
	out := uout.NewOut("JOBMANAGER:")
	// J.runtimeConfig should have Everything in it before now;
	out.V("Runtime Config:")
	out.IndR().ILV(J.runtimeConfig)
	out.H("Job Specs:")
	for _, js := range J.specs {
		out.V(js.briefDetail())
	}
	return out.String()
}
