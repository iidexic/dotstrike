package dscore

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"iidexic.dotstrike/config"
)

func TestOptionNameAndOrder(t *testing.T) {
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
