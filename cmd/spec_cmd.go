/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	pops "iidexic.dotstrike/pathops"
)

type specFlags struct {
	modify, yconfirm bool
}

var flagDataSpec specFlags

// specCmd represents the cfg command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "short descr",
	Long:  `long descr`,
	Run: func(cmd *cobra.Command, args []string) {
		/* ╭───────────────────────── CFG command logic ─────────────────────────╮
		Subcommands: - src - tgt
		//Persistent Flags Short: v,  c, s, t, g(remove),
		} */
		srcarg := *pData.src
		tgtarg := *pData.tgt
		var addSucceed = make(map[string]bool, len(args)+len(*pData.src)+len(srcarg)+len(tgtarg))
		if len(args) > 0 {
			// base args always interpreted as cfg alias
			// will be created if do no exist
			for _, astr := range args {
				found := dscore.SelectSpec(astr)
				// if not found and -y flag or user confirmation, make new cfg
				if !found && (flagDataSpec.yconfirm || askConfirmf("Create new Spec: %s", astr)) {
					td := dscore.GetTempData()
					if td == nil {
						panic(fmt.Errorf("TempData nil"))
					}
					ucfg := td.NewSpec(astr)

					if pData.bsrc {
						var confirmsrc bool
						for _, sa := range srcarg {
							sabasic := pops.MaybeBasicPath(sa) //TODO: Eliminate/Replace
							if !sabasic {
								confirmsrc = askConfirmf(fmt.Sprintf("add non-basic path '%s' to sources?", sa))
							}
							if sabasic || confirmsrc {
								added := ucfg.CheckAddPath(sa, true)
								said := "src:" + sa

								addSucceed[said] = added
							}
						}
						//TODO: FINISH
					}
					if pData.btgt {
						var confirmtgt bool
						for _, ta := range srcarg {
							tabasic := pops.MaybeBasicPath(ta)
							if !tabasic {
								confirmtgt = askConfirmf(fmt.Sprintf("add non-basic path '%s' to targets?", ta))
							}
							if tabasic || confirmtgt {
								added := ucfg.CheckAddPath(ta, false)
								taid := "src:" + ta
								addSucceed[taid] = added
							}
						}
						//TODO: FINISH
					}
					_ = ucfg

				} else {

				}
			}

		}
	},
}

type specOptype int

const (
	_ specOptype = iota
	newSpec
	findSpec
	deleteSpec
	addSrc
	addTgt
	rmvSrc
	rmvTgt
	modOverrides
)

func (f specFlags) identify(pf persistentData) specOptype {
	return 0
type specOpData struct {
	flags         specFlags
	cmd           *cobra.Command
	args          []string
	argcount      int
	argExists     []bool
	existingSpecs []*dscore.Spec
}

func (op *specOpData) populateExisting(args []string) []string {
	notFound := make([]string, 0, len(args))
	for i := range args {
		_ = i
		if spec := dscore.GetSpec(args[i]); spec != nil {
			op.existingSpecs = append(op.existingSpecs, spec)
			op.argExists[i] = true
		} else {
			notFound = append(notFound, args[i])
		}
	}
	return notFound
}
func (op *specOpData) pprintExisting() {
	for i, spec := range op.existingSpecs {
		op.cmd.Print(i, ".   ", spec.Detail())

	}
}

// haveExistingExecute called when have existingSpecs; determines and executes next process steps
func (op *specOpData) haveExistingExecute() {
	if len(op.args) == 1 || len(op.existingSpecs) == op.argcount {
	}
}

func specRun(cmd *cobra.Command, args []string) {
	// if not a spec: make new spec
	// if spec exists, show info
	opSpec := specOpData{
		cmd: cmd, flags: flagDataSpec,
		args: args, argcount: len(args),
		argExists: make([]bool, len(args)),
	}

	notExists := opSpec.populateExisting(args) // may not be necessary to return notExists
	if len(notExists) == len(args) && len(args) > 0 {
		// good to make new ones
	} else if len(opSpec.existingSpecs) > 0 {
		// check what we have
		opSpec.haveExistingExecute()
	}
	_ = notExists

}
func countModFlags() int {
	nlocal := 0
	// if flagDataSpec.modify{
	// 	nlocal++
	// }
	return nlocal + pData.countFlags
}

func (f specOptype) identify(pf persistentData) {
	//flags:
	//persistent:
	//pData
}
func init() {
	rootCmd.AddCommand(specCmd)
	flagDataSpec = specFlags{
		modify:   *specCmd.Flags().BoolP("modify component", "m", false, "modify"),
		yconfirm: *specCmd.Flags().BoolP("autoconfirm user y/n prompts", "y", false, "yes"),
	}

	// PERSISTENT FLAGS:
	//

}
