package dscore

import (
	"fmt"

	pops "iidexic.dotstrike/pathops"
)

//TODO: Make Run command non-persistently modify spec prefs for hard overrides/one-time overrides

func (S *Spec) RunCopy(hardOverride *prefs) error {
	if useHardOverride {

	}
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
			S.applyJobConfig(job)
		}
	}
	return nil
}

func (S *Spec) applyJobConfig(job *pops.CopyJob) *pops.CopyJob {
	var p prefs
	if S.OverrideOn {
		p = S.Overrides
	} else {
		p = gd.data.Prefs
	}

	if p.GlobalTarget {
		job.JobOptionMakeSubdir(true)
	}
	if p.KeepRepo {
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

// For runtime override flag use.
var hardCopyOverride *prefs = &prefs{}
var useHardOverride bool = false

var overrideWhat = map[ConfigOption]bool{
	OptBKeepHidden:   false,
	OptBkeepRepo:     false,
	OptBUseGlobalTgt: false,
}

func MakeHardOverride() *prefs {
	// Add all options to the map, with their existing values in global (or spec if overrides already on)
	useHardOverride = true
	return hardCopyOverride
}

func SetHardOverride(boolOpt ConfigOption, value bool) bool {
	switch boolOpt {
	case OptBkeepRepo:
		hardCopyOverride.KeepRepo = value
	case OptBKeepHidden:
		hardCopyOverride.KeepHidden = value
	case OptBUseGlobalTgt:
		hardCopyOverride.GlobalTarget = value
	}
	return false
}
