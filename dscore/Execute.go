package dscore

import (
	"fmt"
	"slices"

	pops "iidexic.dotstrike/pathops"
)

//TODO: Make Run command non-persistently modify spec prefs for hard overrides/one-time overrides

func BuildCopyJobs() {

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

	if runPrefs.bools[OptBUseGlobalTgt] {
		job.JobOptionMakeSubdir(true)
	}
	if !runPrefs.bools[OptBKeepRepo] {
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
	//*prefs //maybe we just add to map only what we do want to force
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

func SetHardOverride(boolOpt ConfigOption, value bool) bool {
	if slices.Contains(BoolOptions, boolOpt) {
		hardOverrides.options[boolOpt] = value
		return true
	}

	return false
}
