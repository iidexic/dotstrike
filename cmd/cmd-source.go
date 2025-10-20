/*
Copyright © 2025 derek :)
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	"iidexic.dotstrike/uout"
)

type componentCmd struct {
	*cmdData
	alias        *string
	ignore       *[]string
	delete       *bool
	y            *bool
	spec         *[]string
	selectedSpec *bool
	isSource     bool
}

// srcCmd represents the src command
var srcCmd = &cobra.Command{
	Use:   "src path ...",
	Short: "manage sources of a given spec",
	Long:  `add, modify, and delete source components of selected/specified spec`,
	Run:   sourceRun,
}
var src = componentCmd{}

func sourceRun(cmd *cobra.Command, args []string) {
	src.cmdData = newCmdData(cmd, args)
	src.isSource = true
	specFlagArgs := *src.spec
	ns := src.getSpecs(false, specFlagArgs...)

	if len(args) > 0 {
		for _, spec := range src.specs {
			src.components = append(src.components,
				spec.GetMatchingComponents(src.args, src.isSource)...)
		}

	}

	if ns > 0 {
		e := runComponent(&src)
		if e != nil {
			cmd.Print(e.Error())
		}
	} else if ns == 0 {
		cmd.Print("no specs found for entered arguments!")
	}

}

// TODO:(VERY HIGH) Src Deletion tries to create. Target creation makes source and says "failed (already exists)" FIX PLS

func runComponent(cmp *componentCmd) error {
	// TODO: 1 arg w/multiple specs; component doesn't exist in all specs

	numargs, numcomp := len(cmp.args), len(cmp.components)
	if numargs > 0 {

	}
	switch {
	//TODO:(hi) seems like we are NOT handling local  dir correctly. Fix pathComponent functions (should not need to handle on create or delete in cmd)

	// no components found from args -> make new args
	case numargs > 0 && numcomp == 0:
		conftxt := fmt.Sprintf("Add path(s) as %s to multiple specs", componentTypeString(cmp.isSource))
		if len(cmp.specs) > 1 && !checkConfirm(conftxt, cmp.y) {
			return nil
		}
		return cmp.addAll()
	// components found from args but not all args -> uncertain/mixed qty
	case numargs > 0 && numcomp < numargs && numcomp > 0: // Uncertain/Mixed Qty
		cmp.runMixedQty()
	// components found from args and all args (effectively)
	case numcomp > 0:
		switch {
		case *cmp.delete && len(*cmp.ignore) > 0: //Delete Ignores
			if len(cmp.components) == 0 {
				cmp.deleteIgnores()
			}

		case *cmp.delete: //Delete components
			if ls := len(cmp.specs); ls == 1 || *cmp.y ||
				askConfirmf("Delete %d %ss from %d specs?", numcomp, componentTypeString(cmp.isSource), ls) {
				cmp.deleteComponents()
			}
		case len(*cmp.ignore) > 0: // Add Ignores
			return cmp.addIgnores()
		}
	// have specs, no args, no components
	case numcomp == 0 && numargs == 0 && len(cmp.specs) > 0:
		cmp.noArgOutput()
	default:
		cmp.Print("this should never happen")
	}

	return nil
}
func (C componentCmd) noArgOutput() {
	out := uout.NewOutf("--%ss:--", componentTypeString(C.isSource))
	if len(C.specs) == 0 {
		C.Print("No specs found/selected")
	} else {
		for _, spec := range C.specs {
			out.IfF(spec == dscore.TempData().SelectedSpec(), "Spec %s (selected)", "Spec %s", spec.Alias, spec.Alias)
			out.IfV(C.isSource, spec.DetailSources(false), spec.DetailTargets(false))
		}
		C.Print(out.String())
	}
}

func (C *componentCmd) runMixedQty() error {
	switch {
	case len(*C.ignore) > 0:
		if *C.delete {
			return C.deleteIgnores()
		} else {
			return C.addIgnores()
		}
	case *C.delete:
		if askConfirmf("Found %d sources from %d arguments. Delete?", len(C.components), len(C.args)) {
			C.deleteComponents()
		} else {
			C.Println("Cancelled.")
		}
	default:
		// must have >1 spec or >1 arg to reach mixed qty
		if askConfirmf("Some of the provided arg paths exist in specs and some do not. Try to add nonexisting to all?") {
			C.addAll()
		} else {
			C.Println("Cancelled.")
		}

	}
	return nil
}

func (C *componentCmd) addAll() error {
	cmptype := componentTypeString(C.isSource)
	out := uout.NewOutf("-- add %s(s) --", cmptype)
	out.WipeOnOutput(true)
	for _, spec := range C.specs {
		added := spec.CheckAddMultiplePaths(C.args, C.isSource)
		out.F("spec %s:", spec.Alias)
		out.IndR()
		for i := range added {
			out.NV(C.args[i], added[i])
			if !added[i] {
				out.A(" (path exists as source or target in spec)")
			}
		}
		C.Print(out.String())
		out.Clear() // is there a reason to print each spec?
	}
	return nil
}

// shouldn't need this, oh well for now
func componentTypeString(isSource bool) string {
	if isSource {
		return "source"
	} else {
		return "target"
	}
}

func (C *componentCmd) deleteComponents() int {
	temp := dscore.TempData()
	deleted := make([]string, 0, len(C.components))
	for _, cmp := range C.components {
		deleted = append(deleted, cmp.Descriptor())
		e := temp.GetSpec(cmp.Parent).DeleteByPtr(cmp)
		if e != nil {
			C.Printf("Error on delete: %v", e)
		}
	}
	C.Printf("Deleted %d components: %v", len(C.components), deleted)
	return len(C.components)
}

// TODO: Finish all
func (C *componentCmd) addIgnores() error {
	for _, cmp := range C.components {
		n := cmp.Ignores.Add(*C.ignore...)
		if li := len(*C.ignore); n == li {
			C.Printf("%d ignore patterns added to %s", n, cmp.Descriptor())
		} else {
			C.Printf("%d/%d ignore patterns added to %s", n, li, cmp.Descriptor())

		}
	}
	return nil
}

func (C *componentCmd) deleteIgnores() error {
	for _, cmp := range C.components {
		e := cmp.Ignores.Delete(*C.ignore...)
		if e != nil {
			C.Printf("Error on delete: %v", e)
		}
	}
	return nil
}

func makeCmpFlags(cmd *cobra.Command, cmp *componentCmd) {
	cmp.ignore = cmd.Flags().StringArray("ignore", nil, "ignore")
	cmp.alias = cmd.Flags().String("alias", "", "set alias")
	cmp.delete = cmd.Flags().Bool("delete", false, "delete")
	cmp.y = cmd.Flags().BoolP("yes", "y", false, "Auto-confirm on prompt")
	cmp.spec = cmd.Flags().StringSlice("spec", []string{}, `--spec="alias1, alias2"  to target specs provided`)
	cmp.selectedSpec = cmd.Flags().BoolP("selected", "s", true, "--useSelected=false to disable operating on selected spec")
}

func init() {
	rootCmd.AddCommand(srcCmd)
	makeCmpFlags(srcCmd, &src)
}
