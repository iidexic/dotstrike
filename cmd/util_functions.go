package cmd

import (
	"slices"

	"github.com/spf13/cobra"
)

const strAlphaLower = "abcdefghijklmnopqrstuvwxyz"
const strAlphaUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const strEnclosingOpen = `[{(<`
const strEnclosingMulti = "'`" + `"`
const strEnclosingClosed = `]})>`
const strBasicPunctuation = ",.;:-!?'"
const strOtherSymbols = "~@#$%^&*()_=+/[]\\<>?{}|"
const strNumeric = "0123456789"

type runeTypeCount struct {
	alpha, alphaLower, alphaUpper  int
	enclOpen, enclClose, enclMulti int
	punc, symbols, numeric         int
}

func printNumberedListFiltered(cmd *cobra.Command, textlist []string, filter []bool) {
	for i, text := range textlist {
		cmd.Printf(" %d. %s", i, text)
	}
}

// untested, took different route
func stripFalse(list []any, keep []bool) {
	for i := 1; i < len(list)-1; i++ {
		if !keep[i] {
			keep = slices.Concat(keep[:i], keep[i+1:])
		}

	}
}

func isEven(n int) bool { return n%2 == 0 }
