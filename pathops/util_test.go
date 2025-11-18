package pops

import (
	"testing"
)

var rootinTest = `d:\coding\exampleFiles\INPUT`
var inpathsTest = []string{
	`d:\coding\exampleFiles\INPUT\`,
	`d:\coding\exampleFiles\INPUT\cplusplus-source.cc`,
	`d:\coding\exampleFiles\INPUT\file_format`,
	`d:\coding\exampleFiles\INPUT\pics`,
	`d:\coding\exampleFiles\INPUT\file_format\bad-gif.gif`,
	`d:\coding\exampleFiles\INPUT\file_format\icon-circle_davinci_RAW.raw`,
	`d:\coding\exampleFiles\INPUT\file_format\list.gif`,
	`d:\coding\exampleFiles\INPUT\file_format\spinner.gif`,
	`d:\coding\exampleFiles\INPUT\file_format\tai-ku.gif`,
	`d:\coding\exampleFiles\INPUT\file_format\tinyfolder.gif`,
	`d:\coding\exampleFiles\INPUT\pics\1992 bombing conspiracy theories`,
	`d:\coding\exampleFiles\INPUT\pics\1992 bombing conspiracy theories\35-free-drum-loops.jpg`,
	`d:\coding\exampleFiles\INPUT\pics\1992 bombing conspiracy theories\Electronisounds-SomethingForNothing.jpg`,
	`d:\coding\exampleFiles\INPUT\pics\1992 bombing conspiracy theories\mr_9999_brick_game_9999_in_1.jpg`,
	`d:\coding\exampleFiles\INPUT\pics\1992 bombing conspiracy theories\secret_folder_no_files`,
}
var inpathsResult = []string{
	`d:\coding\exampleFiles\OUTPUT\test-addrel\`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\cplusplus-source.cc`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\file_format`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\pics`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\file_format\bad-gif.gif`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\file_format\icon-circle_davinci_RAW.raw`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\file_format\list.gif`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\file_format\spinner.gif`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\file_format\tai-ku.gif`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\file_format\tinyfolder.gif`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\pics\1992 bombing conspiracy theories`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\pics\1992 bombing conspiracy theories\35-free-drum-loops.jpg`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\pics\1992 bombing conspiracy theories\Electronisounds-SomethingForNothing.jpg`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\pics\1992 bombing conspiracy theories\mr_9999_brick_game_9999_in_1.jpg`,
	`d:\coding\exampleFiles\OUTPUT\test-addrel\pics\1992 bombing conspiracy theories\secret_folder_no_files`,
}

func TestEndsWith(t *testing.T) {
	word := "Bringo"
	checkvs := []string{"in", "ing", "o", "ngo", "ringo", "bringo", "Bringo"}
	for _, v := range checkvs {
		t.Logf("%s endswith %s? [%t]", word, v, endswith(word, v))
	}
}

func TestStartsWith(t *testing.T) {
	word := "Bringo"
	checkvs := []string{"B", "ing", "Bri", "ngo", "ringo", "bringo", "Bringo"}
	for _, v := range checkvs {
		t.Logf("%s startswith %s? [%t]", word, v, startswith(word, v))
	}
}

func TestAddRel(t *testing.T) {
	ins := inpathsTest
	outpath := `d:\coding\exampleFiles\OUTPUT\test-addrel\`
	for _, inp := range ins {
		out, e := addRelpath(outpath, rootinTest, inp)
		t.Logf("in: %s, out: %s", inp, out)
		if e != nil {
			t.Errorf("addRelpath error: %v", e)
		}
	}
}

func TestStripRoot(t *testing.T) {
	ins := inpathsTest
	t.Logf("Root: %s", rootinTest)
	t.Logf("cleanroot: %s", CleanPath(rootinTest))
	t.Logf("Inp[0] == Root? %t", CleanPath(ins[0]) == CleanPath(rootinTest))
	//outpath := `d:\coding\exampleFiles\OUTPUT\test-addrel\`
	for _, inp := range ins {
		out, e := stripRoot(rootinTest, inp)
		t.Logf("inpath: %s, ROOT: %s", inp, out)
		if e != nil {
			t.Errorf("stripRoot(%s) error: %v", inp, e)
		}
	}
}
