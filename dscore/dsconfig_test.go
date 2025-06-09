package dscore

import (
	"strings"
	"testing"
)

var testStringML = ` How do we:
		1. split the string at the comment; remove comment
		2. Use [Table.Details] in that format to split to separate objects
		3. Return usable data; in proper format
		`
var testString = "single line test-string: split [1]:[2]; or something like that :)"

func TestStringOps(t *testing.T) {
	t.Logf("(Original String: %s)", testString)
	t.Log("|- STAR -|")
	t.Log("SPLIT ON INDEX: first colon")
	ci := strings.Index(testString, ":")
	t.Logf("Index=%d", ci)
	part1 := testString[:ci]
	part2 := testString[ci:]
	t.Logf("%+v", part1)
	t.Logf("%+v", part2)
	t.Log("|- END -|")
	t.Fail()

}
