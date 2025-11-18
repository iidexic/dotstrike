package pops

import (
	"fmt"

	"iidexic.dotstrike/config"
	"iidexic.dotstrike/uout"
)

type JobGroup struct {
	pathSet
	groupName   string
	jobNames    []string
	jobPtrs     []*CopyJob
	bcfg        boolConfig
	scfg        stringConfig
	initialized bool
}

func (g JobGroup) String() string { return g.Detail() }

func (g *JobGroup) Name() string { return g.groupName }

func (g *JobGroup) Config() map[config.OptionKey]bool { return g.bcfg }

func (g *JobGroup) makeJobs() {
	for x, in := range g.ins {
		for y, out := range g.outs {
			i := x*len(g.outs) + y
			jname := fmt.Sprintf("%s.in-%d.out-%d", g.groupName, x, y)
			g.jobPtrs[i] = cmachine.NewJob(jname, in, out)
			g.jobNames[i] = jname
		}
	}

}

// RunAll runs all jobs in the group. If abortOnError, a job error will abort the remaining jobs.
//
// IF jobs have no config (len(BPrefs) == 0) and the group has config (len(bcfg) > 0),
// the jobs will be configured with the group's config.
func (g *JobGroup) RunAll(abortOnError bool) error {
	var outError error
	for i := range g.jobPtrs {
		if len(g.jobPtrs[i].BPrefs) == 0 && len(g.bcfg) > 0 {
			g.jobPtrs[i].BPrefs = g.bcfg
		}
		e := g.jobPtrs[i].RunFS()
		if e != nil && abortOnError {
			return e
		} else if e != nil {
			if outError == nil {
				outError = fmt.Errorf("Copy Errors: %w", e)
			} else {
				outError = fmt.Errorf("%w\n%w", outError, e)
			}
		}
	}
	return outError
}

func (g *JobGroup) ConfigToJobs() {
	for i := range g.jobPtrs {
		g.jobPtrs[i].BPrefs = g.bcfg
	}
}

func (g *JobGroup) CopyJobs() []*CopyJob { return g.jobPtrs }

// Same as String()
func (g *JobGroup) Detail() string {
	if !g.initialized {
		if g.groupName != "" {
			return fmt.Sprintf("Job Group: %s (uninitialized)", g.groupName)
		} else {
			return "(uninitialized group)"
		}
	}
	out := uout.NewOutf("Job Group: %s", g.groupName)
	if len(g.bcfg) > 0 {
		out.IndR().V("Group Config:")
		out.ILV(g.bcfg)
	}
	if len(g.scfg) > 0 {
		out.IndR().V("String Config:")
		out.IndR().ILV(g.scfg)
	}
	// WARN: Might just print pointer addresses
	if len(g.jobPtrs) > 0 {
		out.H("Jobs:")
		out.ILV(g.jobPtrs)
	}
	return out.String()

	// ──────────────── old Detail ───────────────────────────────────────
	// This one had the nice pipe format
	/*
		sd := make([]string, len(g.jobPtrs)+len(g.bcfg)+len(g.scfg)+3)
		sd[0]=fmt.Sprintf("|[Group:%s]%d jobs",g.groupName,len(g.jobPtrs))
		i:=1; if len(g.jobPtrs) > 0 {; sd[i] = "|--- JOBS:"; i++
			for ijp, j := range g.jobPtrs {; detailText := fmt.Sprintf("|---	[%d] ", ijp) + j.DetailLine(); sd[i] = detailText; i++; } }
		if len(g.scfg)+len(g.bcfg) > 0 {; sd[i] = "|-- GROUP CONFIG:"; i++
			for k, bv := range g.bcfg {; sd[i] = fmt.Sprintf("|---	%s: %t", k.String(), bv); i++; }
			for k, v := range g.scfg {; sd[i] = fmt.Sprintf("|---	%s: '%s'", k.String(), v); }; }
			   sd = slices.Clip(sd); return strings.Join(sd, "\n")
	*/
}
