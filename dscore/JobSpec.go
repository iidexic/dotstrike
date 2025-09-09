package dscore

import (
	"fmt"
	"maps"

	"iidexic.dotstrike/config"
	pops "iidexic.dotstrike/pathops"
)

type jobSpec struct {
	*Spec
	config                  map[config.OptionKey]bool
	group                   *pops.JobGroup
	manualSpec, partialSpec bool
	configApplied, madeJobs bool
}

// TODO:(low) Move to be method of jobSpec
func (S *Spec) jobName(isrc, itgt int) string {
	return fmt.Sprintf("%s.src-%d.tgt-%d", S.Alias, isrc, itgt)
}

func (js *jobSpec) briefDetail() string {
	detail := fmt.Sprintf("spec %s - ", js.Alias)
	if js.partialSpec {
		detail += "partial, "
	}
	if js.manualSpec {
		detail += "manual, "
	}
	detail += fmt.Sprintf("sources:  %d, targets: %d", len(js.Sources), len(js.Targets))

	return detail

}

func (js *jobSpec) applyConfigsPrioritized(lowPriority, highPriority map[ConfigOption]bool) error {
	if js.config == nil {
		mlen := len(lowPriority) + len(highPriority)
		if js.OverrideOn {
			mlen += len(js.Overrides.Bools)
		}
		js.config = make(map[ConfigOption]bool, mlen)
	}
	maps.Copy(js.config, lowPriority)
	if js.OverrideOn && len(js.Overrides.Bools) > 0 {
		maps.Copy(js.config, js.Overrides.Bools)
	}
	maps.Copy(js.config, highPriority)
	return nil
}

// this is seriously only for input into NewJobGroup as that function is a mess
func (js *jobSpec) groupExport() (string, []string, []string, map[ConfigOption]bool) {
	return js.Alias, js.sourcePaths(), js.targetPaths(), js.config
}
func (js *jobSpec) AddGlobalTarget() {
	if globTarget := tempData.GlobalTargetPath; !js.IsPathChild(globTarget) {
		js.addSources(globTarget)
	}
	//TODO:(low) make this function more flexible

}

// TODO: Finish RemoveGlobalTarget
//
//	 Just do a for lookp thru components, dw bout this other shit.
//	But I am also drunk right now so take with a grain o salt
func (js *jobSpec) RemoveGlobalTarget() {
	if globTarget := tempData.GlobalTargetPath; js.IsPathChild(globTarget) {
		if js.IsPathSource(globTarget) {

		}
	}

}
