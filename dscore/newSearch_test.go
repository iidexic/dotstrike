package dscore

import "testing"

var testGD = globalData{
	Specs: []spec{
		{Alias: "vi"}, {Alias: "nvim"}, {Alias: "neovide"}, {Alias: "vim"}, {Alias: "test"}, {Alias: "Wezterm"},
	},
}

func TestFindSpecs(t *testing.T) {
	tsearchpat := []string{"vim", "v", "no", "e", "ozyboy"}
	for i, s := range tsearchpat {
		out, oe := testGD.FFindSpec(s)
		_, _, _ = out, oe, i
	}
}
