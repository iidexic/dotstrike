package config

import (
	"slices"
	"strings"
)

// LookupOption resolves user input to a single OptionKey.
// Exact matches (LookupExacts) win over substring matches.
// Returns NotAnOption if input matches zero opts OR more than one — silent
// misroute would be worse than a lookup failure for a config command.
// Use LookupOptionCandidates to also get the ambiguous candidates.
func LookupOption(input string) OptionKey {
	opt, _ := LookupOptionCandidates(input)
	return opt
}

// LookupOptionCandidates returns the resolved OptionKey plus the ambiguous
// candidate list when resolution fails due to ambiguity.
//   - Exactly one match (exact or substring): (opt, nil)
//   - Zero matches: (NotAnOption, nil)
//   - Multiple matches: (NotAnOption, candidates)
//
// Exact matches take precedence: if any opt has input in its LookupExacts,
// substring hits are ignored.
func LookupOptionCandidates(input string) (OptionKey, []OptionKey) {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return NotAnOption, nil
	}
	var exactHits, subHits []OptionKey
	for id, opt := range AllOptions {
		if slices.Contains(opt.LookupExacts, input) {
			exactHits = append(exactHits, id)
			continue
		}
		if len(opt.LookupSubstrings) == 0 {
			continue
		}
		match := true
		for _, substr := range opt.LookupSubstrings {
			if !lookupSubstringMatch(input, substr) {
				match = false
				break
			}
		}
		if match {
			subHits = append(subHits, id)
		}
	}
	switch {
	case len(exactHits) == 1:
		return exactHits[0], nil
	case len(exactHits) > 1:
		return NotAnOption, exactHits
	case len(subHits) == 1:
		return subHits[0], nil
	case len(subHits) > 1:
		return NotAnOption, subHits
	}
	return NotAnOption, nil
}

func OptByNameExact(optionName string) OptionKey {
	for k, v := range AllOptions {
		if optionName == v.NameText {
			return k
		}
	}
	return NotAnOption
}

func ConfigsMatch(a, b map[OptionKey]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if v != b[k] {
			return false
		}
	}
	return true
}

// ToConfigMap makes a ConfigMap from a map[string]bool
// each key is converted to an OptionKey. Any real option will be copied to returned ConfigMap
func ToConfigMap(opts map[string]bool) ConfigMap {
	cfg := make(ConfigMap, len(opts))
	for k, v := range opts {
		cfg[LookupOption(k)] = v
	}
	return cfg
}

// CopyToConfig will perform all option lookups and copy to cfg
// If force is true, it will overwrite any existing option, otherwise only new options will be copied
// If cfg is nil, a new map will be created
func CopyToConfig(stringopts map[string]bool, cfg ConfigMap, force bool) (ConfigMap, []string) {
	fails := make([]string, 0, len(stringopts))
	if cfg == nil {
		cfg = make(ConfigMap, len(stringopts))
	}
	if force {
		for k, b := range stringopts {
			if opt := LookupOption(k); opt != NotAnOption {
				cfg[opt] = b
			} else {
				fails = append(fails, k)
			}
		}
	} else {
		for k, b := range stringopts {
			opt := LookupOption(k)
			if _, ok := cfg[opt]; opt != NotAnOption && !ok {
				cfg[opt] = b
			} else if opt == NotAnOption {
				fails = append(fails, k)
			}
		}
	}
	return cfg, fails
}
