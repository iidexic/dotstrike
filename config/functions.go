package config

import (
	"slices"
	"strings"
)

// NOTE: Added check of LookupExacts to make life easier
func LookupOption(input string) OptionKey {
	input = strings.TrimSpace(strings.ToLower(input))
	for id, opt := range AllOptions {
		match := true
		for _, substr := range opt.LookupSubstrings {
			match = match && lookupSubstringMatch(input, substr)
		}
		if match || slices.Contains(opt.LookupExacts, input) {
			return id
		}
	}
	return NotAnOption
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
