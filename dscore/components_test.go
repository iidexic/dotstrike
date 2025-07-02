package dscore

import (
	"testing"

	pops "iidexic.dotstrike/pathops"
)

func testcfg() *cfg {
	c := cfg{}

	return &c
}

var cftest cfg = cfg{
	Alias: "tilde_dotconfig",
	Sources: []pathComponent{
		{
			Alias:   "Source1",
			Path:    "~/.config/",
			Abspath: pops.HomeDirtyJoin(".config/"),
		},
	},
}

func TestInherent(t *testing.T) {

}
