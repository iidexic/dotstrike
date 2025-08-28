package cmd

import (
	"github.com/spf13/pflag"
	"iidexic.dotstrike/config"
)

// I keep creating and deleting this file. Just leave it for now

type flagtype = config.ValueType

var (
	tbool    flagtype = config.Tbool
	tstring  flagtype = config.Tstring
	tstrings flagtype = config.TstringSlice
)

type flagdata struct {
	name, short, usage string
	isP, isOption      bool
	datatype           flagtype
	defaultValue       any
	actualValue        any
}

func (f *flagdata) boolDefaultVal() *bool {
	val, ok := f.defaultValue.(*bool)
	if ok {
		return val
	}
	return nil
}
func (f *flagdata) stringDefaultVal() *string {
	val, ok := f.defaultValue.(*string)
	if ok {
		return val
	}
	return nil

}
func (f *flagdata) stringSliceDefaultVal() *[]string {
	val, ok := f.defaultValue.(*[]string)
	if ok {
		return val
	}
	return nil
}

func (f *flagdata) boolVal() *bool { return nil }
func (f *flagdata) stringVal() *string {
	return nil
}
func (f *flagdata) stringSliceVal() *[]string {
	return nil
}

// NOTE: pflags.Flag struct fields:
// Name                string              // name as it appears on command line
// Shorthand           string              // one-letter abbreviated flag
// Usage               string              // help message
// Value               Value               // value as set
// DefValue            string              // default value (as text); for usage message
// Changed             bool                // If the user set the value (or if left to default)
// NoOptDefVal         string              // default value (as text); if the flag is on the command line without any options
// Deprecated          string              // If this flag is deprecated, this string is the new or now thing to use
// Hidden              bool                // used by cobra.Command to allow flags to be hidden from help/usage text
// ShorthandDeprecated string              // If the shorthand of this flag is deprecated, this string is the new or now thing to use
// Annotations         map[string][]string // used by cobra.Command bash autocomple code

var fs = pflag.NewFlagSet("shared", pflag.ExitOnError)
