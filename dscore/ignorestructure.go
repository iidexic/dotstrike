package dscore

import (
	"slices"

	"iidexic.dotstrike/ignore"
	"iidexic.dotstrike/uout"
)

// WARNING: oh god I'm gonna have to do a bunch more toml shit eww
// Actually what if I just keep this shit stored as strings
type (
	ignoreptn  = ignore.TextPattern
	subptn     = ignore.SubPattern
	ignorelist = []ignoreptn
)

type preIgnoreList []string

func (p preIgnoreList) String() string {
	out := uout.NewOut("Ignore Patterns:")
	out.FlatLV(p)
	return out.String()
}

// NOTE: the bool is just for a test/doublecheck get rid of it after
// add adds pattern to the ignore list.
func (p preIgnoreList) add(pattern string) bool {
	l := len(p)
	if !slices.Contains(p, pattern) {
		p = append(p, pattern)
	}
	return len(p) > l
}

// inList just returns slices.Contains(ignoreList, ptn)
func (p preIgnoreList) inList(ptn string) bool { return slices.Contains(p, ptn) }

func (p preIgnoreList) delete(patterns ...string) error {
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
