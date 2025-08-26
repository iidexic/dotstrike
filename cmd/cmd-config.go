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

// configCmd represents the mod command
var configCmd = &cobra.Command{
	Use: "cfg [ optionName optionValue ] ...",
	// ValidArgs - names/altnames for every config to be modified
	ValidArgs: []cobra.Completion{"ignorehidden", "ignorerepo", "useglobaltarget",
		"nofiles", "copyalldirs", "globaltargetpath", "override"},
	Short: "Modify spec overrides and global config values",
	Long: `Change dotstrike configuration options as well as spec config options.
cfg expects arguments in pairs, with each pair containing:
	- Name of the option to be changed
	- New value of that option
To modify global config options, use the --global flag`,
	// global flag as bool? as arg list?
	Run: cfg.run,
}

//TODO: This will just be way easier if it's a struct

type cfgOp struct {
	fYes, fNoSelect, fForceWrite, fGlobal *bool
	fSpec, fSource, fTarget               *[]string
	verbose                               bool
	cmd                                   *cobra.Command
	args                                  []string
	specs                                 []*dscore.Spec
}

func init() {
	rootCmd.AddCommand(configCmd)
	cfg.fYes = configCmd.Flags().BoolP("confirm", "y", false, "--confirm/-y")
	cfg.fNoSelect = configCmd.Flags().Bool("noselect", false, "disable all action on selected spec")
	cfg.fForceWrite = configCmd.Flags().Bool("force-write", false, "force write config to file.")
	cfg.fSpec = configCmd.Flags().StringSlice("spec", []string{}, `--spec="mySpec1,  mySpec2" to target `)
	cfg.fSource = configCmd.Flags().StringSlice("src", []string{}, `--src="srcId1,  c:\path" to select/limit sources`)
	cfg.fTarget = configCmd.Flags().StringSlice("tgt", []string{}, `--tgt="tgtId1,  c:\path" to select/limit targets`)
	cfg.fGlobal = configCmd.Flags().BoolP("global", "g", false, "Operate on global preferences")

}

var cfg cfgOp

func (c *cfgOp) run(cmd *cobra.Command, args []string) {
	la := len(args)
	global := *c.fGlobal
	c.verbose = *pFlags.verbose
	if la == 0 {
		switch {
		case !global && len(*c.fSpec) > 0:
			cfgPrintSelectedSpec(cmd)
			c.cfgPrintFlagSpecs(cmd)
		case *c.fForceWrite:
			dscore.TempData().Modify()
		case global:
			cfgPrintGlobalPrefs(cmd)
		default:
			cfgPrintSelectedSpec(cmd)
		}

	}
	if la >= 2 {
		switch {
		case global:
			c.applyToGlobals(cmd, args)
		case !*c.fNoSelect:
			e := c.applyToSpecs(cmd, args)
			if e != nil {
				cmd.Print(e)
				cfgPrintErrHelp(cmd)
			}
		}
	}
}
func cfgPrintErrHelp(cmd *cobra.Command) {
	cmd.Println("check cfg --help for argument info")
}

func (c *cfgOp) vprintf(s string, vals ...any) {

	if c.verbose {
		if len(vals) == 0 {
			c.cmd.Print(s)
		} else {
			c.cmd.Printf(s, vals...)
		}
	}

}
func (c *cfgOp) vprintSelected() {
	if c.verbose {
		c.cmd.Print("Selected Specs: ")
		for i := range c.specs {
			c.cmd.Printf("%s, ", c.specs[i].Alias)
		}
	}
}

// return == outcome match user intent;
func (c *cfgOp) applyToSpecs(cmd *cobra.Command, args []string) error {
	temp := dscore.TempData()
	specs := getSpecs(cmd, !*c.fNoSelect)
	var confirmUser bool
	if ls := len(specs); ls > 1 {
		confirmUser = checkConfirm(fmt.Sprintf("Apply Options (overrides) to %d specs", ls), c.fYes)
	} else if ls == 1 {
		confirmUser = true
	} else if ls == 0 {
		return fmt.Errorf("no specs selected")
	}
	if confirmUser {
		c.vprintSelected()
		mapargs, remainder := c.cfgArgsMap(args)
		lr := len(remainder)
		c.outputRemainder(remainder)
		if len(mapargs) == 0 {
			//cmd.Print("Failed\n")
			return fmt.Errorf("no config options could be made from args")
		}
		for i := range specs {
			failed := temp.SetSpecOverridesMap(specs[i], mapargs)
			lf := len(failed)
			if lf > 0 {
				cmd.Printf("config options not found for:\n%s", failed)
			}
			switch {
			case lf == 0 && lr == 0:
				cmd.Print("Succesfully wrote all values")
			case lf*2+lr < len(args):
				cmd.Print("Succesfully wrote (with failures)")
			}
		}
	}
	return nil
}

func (c *cfgOp) outputRemainder(remainder []string) {
	lr := len(remainder)
	if !isEven(lr) {
		c.cmd.Printf("unpaired key '%s' not used", remainder[lr-1])
	}
	if lr > 2 {
		c.cmd.Println("cannot make true/false for:")
		for i := range lr / 2 {
			c.cmd.Printf("%s = '%s'", remainder[i*2], remainder[i*2+1])
		}
	}
}

// cfgArgsMap creates a map out of an argument list that can be used in SetSpecOverridesMap.
// Returns the created map, and a string slice containing any values that could not be added to the map.
//
// It assumes that args is formatted as a linear slice of string:"bool" key value pairs
//
//	i.e. []string{"opt1","true","opt2","false"} -> map[string]bool{"opt1":true,"opt2":false}
//
// If a value cannot be converted to a bool, both key/value arg will be added to the remainder.
// If the final key has no corresponding value (when len(args) is odd), that key will also be added to remainder.
func (c *cfgOp) cfgArgsMap(args []string) (map[string]bool, []string) {
	remainder := make([]string, 0, len(args))
	M := make(map[string]bool, len(args)/2)
	_ = M
	for i := 0; i < len(args)-1; i += 2 {
		btry := dscore.StringToBool(args[i+1])
		if btry != nil {
			M[args[i]] = *btry
		} else {
			remainder = append(remainder, args[i], args[i+1])
		}
	}
	return M, remainder
}
func (c *cfgOp) applyToGlobals(cmd *cobra.Command, args []string) {
	temp := dscore.TempData()
	for i := 0; i < len(args)-1; i += 2 {
		opt := dscore.OptionID(args[i])
		switch {
		case dscore.OptionIsBool(opt):
			barg := dscore.StringToBool(args[i+1])
			if barg != nil {
				output := textOptionModified(opt.String(), temp.SetOptionBool(opt, *barg))
				cmd.Print(output)
			} else {
				cmd.Printf("Failed. Cannot convert '%s' to true/false.", args[i+1])
			}

		case dscore.OptionIsString(opt): // unnecessarily messy
			c.cfgApplyGlobalTargetCautious(cmd, args[i+1])
		}
	}
}

func (c *cfgOp) cfgApplyGlobalTargetCautious(cmd *cobra.Command, newpath string) {
	y := false
	temp := dscore.TempData()
	exist, e := pops.PathExists(newpath)
	if e != nil {
		y = checkConfirm("Error checking path. Set as Global Target path anyway", c.fYes)
	} else if !exist {
		y = checkConfirm("Path does not exist or was not found. Set as Global Target path anyway", c.fYes)
	} else {
		y = true
	}
	if y {
		e = temp.SetOptionString(dscore.StringGlobalTargetPath, newpath)
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
func (c *cfgOp) cfgPrintFlagSpecs(cmd *cobra.Command) {
	if len(c.specs) > 0 {
		for _, arg := range c.args {
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

// // what was this for (unused?)
// func makeValid(argcomponents [][]string) []string {
// 	outsl := make([]string, 0, len(argcomponents)*3)
// 	for i := range argcomponents {
// 		outsl = append(outsl, strings.Join(argcomponents[i], ""))
// 		outsl = append(outsl, strings.Join(argcomponents[i], "-"))
// 		outsl = append(outsl, strings.Join(argcomponents[i], "_"))
// 	}
//
// 	return outsl
// }
