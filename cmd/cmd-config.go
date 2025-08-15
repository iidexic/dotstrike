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
	// WARNING: Probably need to delete?
	ValidArgs: []cobra.Completion{"true", "false",
		"KeepHidden", "KeepRepo", "UseGlobalTarget",
		"CopyFiles", "CopyAllDirs", "GlobalTargetPath",
		"keephidden", "keeprepo", "useglobaltarget",
		"copyfiles", "copyalldirs", "globaltargetpath"},
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

func cfgSpecApply(cmd *cobra.Command, args []string) {
	temp := dscore.TempData()
	specs := getSpecs(cmd, !*cfgFlagNoSelect)
	var confirmUser bool
	if ls := len(specs); ls > 1 {
		confirmUser = checkConfirm(fmt.Sprintf("Apply Overrides to %d specs", ls), cfgFlagY)
	} else if ls == 1 {
		confirmUser = true
	}
	if confirmUser {
		for i := range specs {
			e := temp.SetSpecOverridesMap(specs[i], cfgArgsMap(args))
			if e != nil {

			}
		}
	}
}

// cfgArgsMap creates a map out of an argument list that can be used in SetSpecOverridesMap.
//
// It assumes that args is formatted as a linear slice of string:"bool" key value pairs
//   - even index values (args[i]) are treated as a string key (representing a ConfigOption)
//   - the next index value ( args[i+1]) will be used as the bool value for that key.
//
// If a string cannot be transformed into a bool true/false, the key/value will be discarded.
func cfgArgsMap(args []string) map[string]bool {
	M := make(map[string]bool, len(args)/2)
	_ = M
	for i := 0; i < len(args)-1; i += 2 {
		btry := dscore.StringToBool(args[i+1])
		if btry != nil {
			M[args[i]] = *btry
		}
	}
	return M
}

func cfgGlobalApply(cmd *cobra.Command, args []string) {
	temp := dscore.TempData()
	for i := 0; i < len(args)-1; i += 2 {
		opt := dscore.OptionID(args[i])
		switch {
		case dscore.OptionIsBool(opt):
			barg := dscore.StringToBool(args[i+1])
			if barg != nil {
				output := textOptionModified(opt.Text(), temp.SetOptionBool(opt, *barg))
				cmd.Print(output)
			} else {
				cmd.Printf("Failed. Cannot convert '%s' to true/false.", args[i+1])
			}

		case dscore.OptionIsString(opt): // unnecessarily messy
			cfgApplyGlobalTargetCautious(cmd, args[i+1])
		}
	}
}

func cfgApplyGlobalTargetCautious(cmd *cobra.Command, newpath string) {
	y := false
	temp := dscore.TempData()
	exist, e := pops.PathExists(newpath)
	if e != nil {
		y = checkConfirm("Error checking path. Set as Global Target path anyway", cfgFlagY)
	} else if !exist {
		y = checkConfirm("Path does not exist or was not found. Set as Global Target path anyway", cfgFlagY)
	} else {
		y = true
	}
	if y {
		e = temp.SetOptionString(dscore.OptSGlobalTargetPath, newpath)
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

var cfgFlagY, cfgFlagNoSelect *bool

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
	cfgFlagY = configCmd.Flags().BoolP("confirm", "y", false, "--confirm/-y")
	cfgFlagNoSelect = configCmd.Flags().Bool("noselect", false, "disable all action on selected spec")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
