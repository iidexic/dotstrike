/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// flagstorage represents the flagstorage command
var flagstorageCmd = &cobra.Command{
	Use:   "flagstorage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flagstorage called")
	},
}

func init() {
	rootCmd.AddCommand(flagstorageCmd)
}

var errNoCommand error = fmt.Errorf("No cobra.Command attached; need somewhere to put the flag!")

type makeflag[T flagable] interface {
	attach(cmd *cobra.Command)
	make() error
	get() T
}

type flagable interface {
	bool | int | float64 | string | []string
}

func NewFlag[T flagable](T) *anyFlag[T] { return &anyFlag[T]{} }

type flagbase struct {
	cmd            *cobra.Command
	name, usage, p string
	isP            bool
}

// TODO: (loww) figure out how to do this better without needing to have a bunch of typed pointers
type anyFlag[T flagable] struct {
	flagbase
	val            *T
	initVal        T
	boolOut        *bool
	intOut         *int
	floatOut       *float64
	stringOut      *string
	stringSliceOut *[]string
}

func (af *anyFlag[T]) connect(cmd *cobra.Command) {
	af.cmd = cmd
}

func (af *anyFlag[T]) make() error {
	if af.cmd == nil {
		return errNoCommand
	}
	switch any(af.initVal).(type) {
	case bool:
		af.boolOut = af.makeBool()
		return nil
	case int:

	case float64:

	case string:

	case []string:

	}
	return fmt.Errorf("error should be unreachable :)")
}

func (af *anyFlag[T]) makeBool() *bool {
	initb := false
	if af.isP {
		af.boolOut = af.cmd.Flags().BoolP(af.name, af.p, initb, af.usage)
		return af.boolOut
	}
	af.boolOut = af.cmd.Flags().Bool(af.name, initb, af.usage)
	return af.boolOut
}

// func (af *anyFlag[T]) asBool() (*bool, error) {
// 	initb, ok := any(af.initVal).(bool)
// 	if !ok {
// 		return nil, fmt.Errorf("flag %s is not a bool!", af.name)
// 	}
//
// 	if af.isP {
// 		af.boolOut = af.cmd.Flags().BoolP(af.name, af.p, initb, af.usage)
// 		return af.boolOut, nil
// 	}
// 	af.boolOut = af.cmd.Flags().Bool(af.name, initb, af.usage)
// 	return af.boolOut, nil
// }
