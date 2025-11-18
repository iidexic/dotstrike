package dscore

import (
	"slices"

	"iidexic.dotstrike/match"
	"iidexic.dotstrike/uout"
)

type (
	ignoreptn  = match.TextPattern
	subptn     = match.SubPattern
	ignorelist = []ignoreptn
)

type preIgnoreList []string

func (p preIgnoreList) String() string {
	out := uout.NewOut("Ignore Patterns:")
	out.FlatLV(p)
	return out.String()
}

// NOTE: the bool is just for a test/doublecheck get rid of it after
// Add adds pattern to the ignore list.
func (p preIgnoreList) Add(pattern ...string) int {
	l := len(p)
	for _, ptn := range pattern {
		ptn = QuickClean(ptn)
		if !slices.Contains(p, ptn) {
			p = append(p, ptn)
		}
	}
	return len(p) - l
}

// inList just returns slices.Contains(ignoreList, ptn)
// func (p preIgnoreList) inList(ptn string) bool { return slices.Contains(p, ptn) }
func (p preIgnoreList) Delete(patterns ...string) error {
	for _, s := range patterns {
		i := slices.Index(p, s)
		switch {
		case i < 0:
			continue
		case i == 0:
			p = p[1:]
		case i == len(p)-1:
			p = p[:len(p)-1]
		default:
			p = slices.Concat(p[:i], p[i+1:])
		}
	}
	return nil
}
