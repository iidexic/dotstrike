/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"slices"

	"github.com/spf13/cobra"
)

// Wrote all of this just to make flag calls in cmd init() a few characters shorter

// UNUSED
type flaginput struct {
	spec, src, tgt, ignore *[]string
}

// UNUSED

// Laziest shared flag
func addFlagStrSlice(cmd *cobra.Command, flagname string) *[]string {
	if flagname == "spec" {
		return cmd.Flags().StringSlice("spec", []string{}, `--spec="alias1,alias2"`)
	}
	if flagname == "src" {
		return cmd.Flags().StringSlice("src", []string{}, `--src="pathOrAlias,path2"`)
	}
	if flagname == "tgt" {
		return cmd.Flags().StringSlice("tgt", []string{}, `--tgt="pathOrAlias,path2"`)
	}
	if flagname == "ignore" {
		return cmd.Flags().StringSlice("ignore", []string{}, `--ignore="ignoreptn,ignoreptn2"`)
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────
// ──────────────────────────────────────────────────────────────────────
// ──────────── 3rd times a charm ───────────────────────────────────────

var sharedflags *flaglib

func initializeSharedFlags() {
	f := flaglib{}
	f.all = make(map[string]lazyflag)
	f.all["spec"] = &stringSliceFlag{fdata: fdata{name: "spec", usage: `--spec="alias1,alias2" to operate on specs given`}}
	f.all["src"] = &stringSliceFlag{fdata: fdata{name: "src", usage: `--src="id1, id2" to operte on/use sources given`}}
	f.all["tgt"] = &stringSliceFlag{fdata: fdata{name: "tgt", usage: `--tgt="id1, id2" to operte on/use sources given`}}
	f.all["all"] = &boolFlag{fdata: fdata{name: "all", short: "a", usage: "", useShorthand: true}}
	f.all["confirm"] = &boolFlag{fdata: fdata{name: "confirm", usage: "auto-confirm all user prompts", short: "y", useShorthand: true}}
	sharedflags = &f
}

type lazyflag interface {
	attach(*cobra.Command)
	setUsage(string)
	fname() string
	value() any
}

type flaglib struct {
	all            map[string]lazyflag
	boolVal        map[string]*bool
	stringVal      map[string]*string
	stringSliceVal map[string]*[]string
}

type fdata struct {
	name, usage, short string
	useShorthand       bool
}

func (F *flaglib) add(cmd *cobra.Command, names ...string) {
	for n, flag := range F.all {
		if slices.Contains(names, n) {
			flag.attach(cmd)
			switch flag.value().(type) {
			case *bool:
				F.boolVal[flag.fname()] = flag.value().(*bool)
			case *string:
				F.stringVal[flag.fname()] = flag.value().(*string)
			case *[]string:
				F.stringSliceVal[flag.fname()] = flag.value().(*[]string)
			}
		}
	}
}

func (f *fdata) setUsage(usage string) {
	f.usage = usage
}

func (f *fdata) fname() string { return f.name }

type boolFlag struct {
	fdata
	val *bool
}

type stringFlag struct {
	fdata
	val *string
}

type stringSliceFlag struct {
	fdata
	val *[]string
}

func (F *flaglib) setUsage(name, usage string) {
	_, exist := F.all[name]
	if exist {
		F.all[name].setUsage(usage)
	}
}

func (f *boolFlag) attach(cmd *cobra.Command) {
	if f.useShorthand {
		f.val = cmd.Flags().BoolP(f.name, f.short, false, f.usage)
	} else {
		f.val = cmd.Flags().Bool(f.name, false, f.usage)
	}
}
func (f *boolFlag) value() any { return f.val }

func (f *stringFlag) attach(cmd *cobra.Command) {
	if f.useShorthand {
		f.val = cmd.Flags().StringP(f.name, f.short, "", f.usage)
	} else {
		f.val = cmd.Flags().String(f.name, "", f.usage)
	}
}

func (f *stringFlag) value() any { return f.val }

func (f *stringSliceFlag) attach(cmd *cobra.Command) {
	if f.useShorthand {
		f.val = cmd.Flags().StringSliceP(f.name, f.short, []string{}, f.usage)
	} else {
		f.val = cmd.Flags().StringSlice(f.name, []string{}, f.usage)
	}
}

func (f *stringSliceFlag) value() any { return f.val }
