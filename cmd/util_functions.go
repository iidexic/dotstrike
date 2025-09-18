package cmd

import (
	"fmt"
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

type multErr struct {
	error
}

func (m *multErr) Add(errs ...error) {
	for _, e := range errs {
		m.error = fmt.Errorf("%w, %w", m.error, e)
	}
}

func NewMultiError(msg string) multErr { return multErr{error: fmt.Errorf(msg)} }

// func mErr(op string, errs ...error) error {
// 	if len(errs) <= 1 {
// 		return errs[0]
// 	}
// 	var out error
// 	if errs[0] == nil {
// 		out = fmt.Errorf("errs during %s: ", op)
// 		errs = errs[1:]
// 	}
// 	return out
// }

func printNumberedListFiltered(cmd *cobra.Command, textlist []string, filter []bool) {
	for i, text := range textlist {
		cmd.Printf(" %d. %s", i, text)
	}
}

// untested, took different route
// WARNING: This probably bad
func stripFalse(list []any, keep []bool) {
	for i := 1; i < len(list)-1; i++ {
		if !keep[i] {
			keep = slices.Concat(keep[:i], keep[i+1:])
		}

	}
}

func sliceUniques[E comparable](in []E) []E {
	seen := make(map[E]struct{}, len(in))
	out := make([]E, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func listOut[S ~[]E, E any](paren bool, lvals ...E) string {
	sout := ""
	if paren {
		sout = "( "
	}
	if llv := len(lvals); llv > 0 {
		sout = fmt.Sprintf("%v", lvals[0])
		if llv == 1 {

			return sout
		}
		lvals = lvals[1:]
	}

	for _, x := range lvals {
		sout = fmt.Sprintf("%s, %v", sout, x)
	}

	return sout
}

func isEven(n int) bool { return n%2 == 0 }
