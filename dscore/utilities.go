package dscore

import (
	"fmt"
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

// ValuesInSlice returns a unique set of values that are in both s and vals
func ValuesInSlice[E comparable](s []E, vals ...E) []E {
	sSet := make(map[E]struct{}, len(vals))
	for _, v := range vals {
		sSet[v] = struct{}{}
	}

	r := make([]E, 0, min(len(s), len(vals)))

	for _, v := range s {
		if _, ok := sSet[v]; ok {
			delete(sSet, v)
			r = append(r, v)
		}
	}
	return r
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

// NoneIn reports whether s contains anything in vals
func NoneIn[E comparable](s []E, vals ...E) bool {
	for _, val := range vals {
		if slices.Contains(s, val) {
			return false
		}
	}
	return true
}

func ExtendErr(errs ...error) error {
	var eout error
	if le := len(errs); le == 0 {
		return nil
	} else if le == 1 && errs[0] != nil {
		return errs[0]
	}
	eout = errs[0]
	errs = errs[1:]
	for _, e := range errs {
		if e != nil {
			eout = fmt.Errorf("%w, %w", eout, e)
		}
	}
	return eout
}

// func lastCharNumber(s string) (int, bool) {
// 	n := len(s) - 1
// 	if unicode.IsDigit(rune(s[n])) {
// 		i, e := strconv.Atoi(string(s))
// 		if e == nil {
// 			return i, true
// 		}
// 	}
// 	return -1, false
// }

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
