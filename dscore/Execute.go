package dscore

import (
	"fmt"

	"iidexic.dotstrike/config"
	pops "iidexic.dotstrike/pathops"
)

//TODO: Make Run command non-persistently modify spec prefs for hard overrides/one-time overrides

type jobProcessor struct {
	specs         []*Spec
	runtimeConfig map[config.OptionKey]bool
}

type jobSpec struct {
	*Spec
	jobs []*pops.CopyJob
}

var manager = jobProcessor{
	specs:         make([]*Spec, 0, 3), //arbitrary. Add init function if want to use userdata spec qty
	runtimeConfig: make(map[config.OptionKey]bool),
}

func JobManager() *jobProcessor { return &manager }

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

func (S *Spec) RunCopy(hardOverride *prefs) error {
	//NOTE: Where do jobconfig? within makeCopyJobs?
	if !S.allInitialized() {
		return fmt.Errorf("spec not initialized: %s", S.Alias)
	}

	return nil
}

// makeCopyJobs creates all CopyJobs for a source.
// makeCopyJobs assumes any overrides have already been applied to S
func (S *Spec) makeCopyJobs() error {
	for x := range S.Sources {
		for y := range S.Targets {
			job := S.newCopyJob(x, y)
			jobprefs := S.makeJobConfig()
			if jobprefs != nil {

			}
			S.applyJobConfig(job)
		}
	}
	return nil
}

func (S *Spec) makeJobConfig() *prefs {

	return nil
}

func (S *Spec) applyJobConfig(job *pops.CopyJob) *pops.CopyJob {
	var runPrefs prefs
	if S.OverrideOn {
		runPrefs = S.Overrides
	} else {
		runPrefs = gd.data.Prefs
	}

	if runPrefs.Bools[BoolUseGlobalTarget] {
		job.JobOptionMakeSubdir(true)
	}
	if runPrefs.Bools[BoolIgnoreRepo] {
		job.IgnoreGit()
	}

	return job
}

func (S *Spec) newCopyJob(isrc, itgt int) *pops.CopyJob {
	return pops.Copier.NewJob(
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
