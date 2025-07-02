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

type cfgFlags struct {
	modify, yconfirm, all *bool
}

var flagDataCfg cfgFlags

// cfgCmd represents the cfg command
var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "short descr",
	Long:  `long descr`,
	Run: func(cmd *cobra.Command, args []string) {
		/* ╭───────────────────────── CFG command logic ─────────────────────────╮
		Subcommands: - src - tgt
		//Persistent Flags Short: v,  c, s, t, g(remove),
		} */
		srcarg := *pData.src
		tgtarg := *pData.tgt
		fsrc := len(srcarg) > 0
		ftgt := len(tgtarg) > 0
		var addSucceed = make(map[string]bool, len(args)+len(*pData.src)+len(srcarg)+len(tgtarg))
		if len(args) > 0 {
			// base args always interpreted as cfg alias
			// will be created if do no exist
			for _, astr := range args {
				found := dscore.SelectCfg(astr)
				// if not found and -y flag or user confirmation, make new cfg
				if !found && (*flagDataCfg.yconfirm || askConfirmf("Create new cfg: %s", astr)) {
					dscore.InitTempData()
					td := dscore.GetTempData()
					if td == nil {
					}
					ucfg := td.NewCfg(astr)

					if fsrc {
						var confirmsrc bool
						for _, sa := range srcarg {
							sabasic := pops.IsBasicPath(sa)
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
					if ftgt {
						var confirmtgt bool
						for _, ta := range srcarg {
							tabasic := pops.IsBasicPath(ta)
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

func init() {
	rootCmd.AddCommand(cfgCmd)
	flagDataCfg = cfgFlags{
		modify:   cfgCmd.Flags().BoolP("modify component", "m", false, "modify"),
		all:      cfgCmd.Flags().BoolP("apply to all found", "a", false, "all"),
		yconfirm: cfgCmd.Flags().BoolP("autoconfirm user y/n prompts", "y", false, "yes"),
	}

	// PERSISTENT FLAGS:
	//

}
