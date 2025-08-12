/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	pops "iidexic.dotstrike/pathops"
)

var argcmp = [][]string{{"global", "target", "path"}, {"use", "global", "target"},
	{"enable", "global", "target"}, {"keep", "repo"}, {"keep", "git"}, {"copy", "repo"}}

// configCmd represents the mod command
var configCmd = &cobra.Command{
	Use: "cfg [ optionName optionValue ] ...",
	// ValidArgs - names/altnames for every config to be modified
	ValidArgs: []cobra.Completion{"globaltarget", "globaltargetpath", "path", "useGlobalTarget"},

	Short: "Modify spec overrides and global config values",
	Long: `Change dotstrike configuration options as well as spec config options.
cfg expects arguments in pairs, with each pair containing:
	- Name of the option to be changed
	- New value of that option
To modify global config options, use the --global flag`,
	// global flag as bool? as arg list?
	Run: cfgRun,
}

func cfgRun(cmd *cobra.Command, args []string) {
	la := len(args)
	switch {
	case la == 0 && !*pFlags.global:
		cfgPrintSelectedSpec(cmd)
		cfgPrintFlagSpecs(cmd)

	case la == 0:
		cfgPrintGlobalPrefs(cmd)

	case *pFlags.global && len(args) >= 2: //apply args to global config
	case len(args) >= 2:
	}
}

func cfgGlobalApply(cmd *cobra.Command, args []string) {
	temp := dscore.TempData()
	for i := 0; i < len(args)-1; i += 2 {
		opt := dscore.OptionID(args[i])
		switch opt {
		case dscore.OptBoolKeepRepo, dscore.OptBoolKeepHidden, dscore.OptBoolUseGlobalTarget:
			barg := dscore.StringToBool(args[i+1])
			if barg != nil {
				output := textOptionModified(opt.Text(), temp.SetOptionBool(opt, *barg))
				cmd.Print(output)
			} else {
				cmd.Printf("Failed. Cannot convert '%s' to true/false.", args[i+1])
			}

		case dscore.OptStringGlobalTargetPath: // unnecessarily messy
			cfgApplyGlobalTargetCautious(cmd, args[i+1])
		}
	}
}

func cfgApplyGlobalTargetCautious(cmd *cobra.Command, newpath string) {
	y := false
	temp := dscore.TempData()
	exist, e := pops.PathExists(newpath)
	if e != nil {
		y = checkConfirm("Error checking path. Set as Global Target path anyway", flagConfirm)
	} else if !exist {
		y = checkConfirm("Path does not exist or was not found. Set as Global Target path anyway", flagConfirm)
	} else {
		y = true
	}
	if y {
		e = temp.SetOptionString(dscore.OptStringGlobalTargetPath, newpath)
		if e != nil {
			cmd.Printf("Error converting path '%s' to absolute path", newpath)
		}
	}

}

func cfgPrintGlobalPrefs(cmd *cobra.Command) {
	temp := dscore.TempData()
	cmd.Printf(`Global Config Options:
	GlobalTarget Path = '%s'
`, temp.GlobalTargetPath)
	cmd.Print(temp.Prefs.Detail())
}
func cfgPrintSelectedSpec(cmd *cobra.Command) {
	temp := dscore.TempData()
	spec := temp.SelectedSpec()
	if spec != nil {
		cmd.Printf("Spec %s Override Options:\n	Override Enabled: %t\n", spec.Alias, spec.OverrideOn)
		cmd.Print(spec.Overrides.Detail())
	} else {
		warnNilSelectedSpec(cmd)
	}
}

// cfgPrintFlagSpecs finds specs from aliases passed via the spec persistent flag, and outputs their override information.
func cfgPrintFlagSpecs(cmd *cobra.Command) {
	if pFlags.bspec {
		for _, arg := range *pFlags.spec {
			if spec := temp.GetSpec(arg); spec != nil {
				cmd.Printf("Spec %s Override Options:\n	Override Enabled: %t\n", spec.Alias, spec.OverrideOn)
				cmd.Print(spec.Overrides.Detail())
			} else {
				cmd.Printf("No spec %s found\n", arg)
			}
		}
	}

}

// textOptionModified
func textOptionModified(val string, modified bool) string {
	if modified {
		return fmt.Sprintf("succesfully updated %s", val)
	} else {
		return fmt.Sprintf("failed to update %s", val)
	}

}

var flagConfirm *bool

type cfgOpt int

func makeValid(argcomponents [][]string) []string {
	outsl := make([]string, 0, len(argcomponents)*3)
	for i := range argcomponents {
		outsl = append(outsl, strings.Join(argcomponents[i], ""))
		outsl = append(outsl, strings.Join(argcomponents[i], "-"))
		outsl = append(outsl, strings.Join(argcomponents[i], "_"))
	}

	return outsl
}

func init() {
	rootCmd.AddCommand(configCmd)
	flagConfirm = configCmd.Flags().BoolP("confirm", "y", false, "--confirm/-y")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
