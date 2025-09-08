package dscore

import (
	"slices"
)

//NOTE: Cobra already does this inherently for some flags. Figure out if it can be used for args

// StringToBool tries to get a bool from string
// does not use strconv or pflag builtin, most likely use one of those
// if succeeds, returns found value (*true or *false)
// if fails to match with any option, returns nil
func StringToBool(text string) *bool {
	var t bool = true
	var f bool = false
	text = QuickClean(text)
	switch text {
	case "true", "1", "yes", "t", "y", "on", "enabled":
		return &t
	case "false", "0", "no", "f", "n", "off", "disabled":
		return &f
	default:
		return nil

	}
}

func KeepIndices[A any](s []A, ikeep []int) []A {
	if len(s) == 0 || len(ikeep) == 0 {
		return []A{}
	}
	if len(ikeep) == 1 && ikeep[0] < len(s) {
		i := ikeep[0]
		return s[i:i]
	}
	out := make([]A, len(ikeep))
	slices.Sort(ikeep)
	offset := 0
	for i, n := range ikeep {
		if i > 0 && n == ikeep[i-1] {
			offset++ //prevent gaps
			continue
		}
		out[i-offset] = s[n]
	}
	out = out[0 : len(ikeep)-offset]
	return out
}

// TODO: clean up; no use for these I can think of.
//
// // StringToBoolFalsy returns true only if text == "true" (case insensitive, spaces removed)
// // returns false in any other case
// func StringBoolTrueOnly(text string) bool {
// 	text = strings.TrimSpace(strings.ToLower(text))
// 	if text == "true" {
// 		return true
// 	}
// 	return false
// }
//
// func StringBoolTruthyFalsy(text string) bool { return len(text) > 0 }
