package dscore

// StringToBool tries to get a bool from string
// if succeeds, returns found value (*true or *false)
// if fails to match with any option, returns nil
func StringToBool(text string) *bool {
	var t bool = true
	var f bool = false
	text = quickclean(text)
	switch text {
	case "true", "1", "yes", "t", "y", "on", "enabled":
		return &t
	case "false", "0", "no", "f", "n", "off", "disabled":
		return &f
	default:
		return nil

	}
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
