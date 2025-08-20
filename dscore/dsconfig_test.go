package dscore

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"iidexic.dotstrike/config"
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

func TestConfigOptionOrder(t *testing.T) {
	countFailures := 0

	for i, o := range config.AllOptionIDs() {
		indexFail := false
		icfgopt := ConfigOption(i)
		itext := icfgopt.String()
		if icfgopt != o {
			t.Errorf("ConfigOption(%d) = %s, should equal %s", i, itext, o.String())
			indexFail = true
		}
		l := OptionID(strings.ToLower(itext))
		if l != icfgopt {
			t.Errorf("Get ID: Input [%d]'%s' - Got [%d]'%s' ", i, itext, int(l), l.String())
			indexFail = true
		}

		if indexFail {
			countFailures++
		}
	}
	t.Logf("Checked ConfigOption System: %d failures.", countFailures)
}

func TestMarshal(t *testing.T) {
	CoreConfig()
	t.Logf("%v", gd.data.Prefs.Bools)
	data, e := toml.Marshal(gd.data.Prefs)
	if e != nil {
		t.Errorf("Error:%s", e.Error())
	}
	t.Log(string(data))
}
