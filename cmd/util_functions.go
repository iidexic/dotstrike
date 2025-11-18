package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func printNumberedListFiltered(cmd *cobra.Command, textlist []string, filter []bool) {
	for i, text := range textlist {
		if filter[i] {
			cmd.Printf(" %d. %s", i, text)
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
