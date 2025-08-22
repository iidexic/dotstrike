package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"text/scanner"

	"iidexic.dotstrike/dscore"
)

var lesstext = "what huhh??!?!1?"
var sometext = `What is happening here?
Why should you care? Should you care?
...
I will answer that question with 2 more:first; what is caring? second; are u dumb? lol
also here's some other symbols !@#$~<>~ and whatever this is `
var stringslice = []string{"1", "two", "[3]", "i!=4", lesstext}

func testingStringSlice() []string {
	p1 := "Header Text:"
	p2 := "-----------"
	p3 := " * This comes next"
	p4 := " * And then this one"
	p5 := "okay, all done"
	outtext := make([]string, 0, 5)
	outtext = append(outtext, p1, p2, p3, p4, p5)
	return outtext
}

func getptrOrClosest(s []string, i int) *string {
	if len(s) > i {
		return &s[i]
	}
	return &s[len(s)-1]
}
func getcwd(t *testing.T) string {
	cwd, e := os.Getwd()
	if e != nil {
		t.Logf("[cwd: %e]", e)
	}
	return cwd
}

func retAsIs(fpath string) error {
	f, e := os.OpenFile(fpath, os.O_RDONLY, 0)
	f.Close()
	return e
}
func retCheck(fpath string) error {
	f, e := os.OpenFile(fpath, os.O_RDONLY, 0)
	f.Close()
	if e != nil {
		return e
	}
	return nil
}

func TestRangeEnd(t *testing.T) {
	t.Log("Range 10:")
	s := ""
	for i := range 10 {
		s += fmt.Sprintf(" %d", i)
		if i == 10 {
			print("Hit 10")
		}
	}
	t.Log(s)
}

func TestPrintEnum(t *testing.T) {
	t.Log("printing dscore.OptBoolKeepRepo (v, +v, #v, s maybe):")
	t.Logf("[%v]", dscore.BoolIgnoreRepo)
	t.Logf("[%+v]", dscore.BoolIgnoreRepo)
	t.Logf("[%#v]", dscore.BoolIgnoreRepo)
	t.Logf("[%#+v]", dscore.BoolIgnoreRepo)
	t.Logf("[%x]", dscore.BoolIgnoreRepo)
}

// func TestIndexBrainpower(t *testing.T) {
// 	si := []int{0, 1, 2, 3, 4, 5, 6, 7}
// 	siw := []string{"big", "stuff", "we", "are", "doing", "things", "DONT PRINT THIS"}
// 	for i := 0; i < len(si)-1; i += 2 {
// 		t.Logf("ints:[%d:%d]", si[i], si[i+1])
// 	}
// 	for i := 0; i < len(siw)-1; i += 2 {
// 		t.Logf("strings:[%s:%s]", siw[i], siw[i+1])
// 	}
// }

func TestIndexString(t *testing.T) {
	text := lesstext
	for i := range text {
		t.Logf("[%d] - %s", i, string(text[i]))
	}
}

func TestScan(t *testing.T) {
	// note 1: scanner.TokenString() takes rune x and returns string "x"
	t.Log("Initial text: ", lesstext)
	nuber, er := fmt.Scan(lesstext)
	_, _ = nuber, er
	_ = scanner.Char
}
func TestTypeOf(t *testing.T) {
	ss := testingStringSlice()
	ss0 := ss[0:0]
	ssp := &ss
	ssp0 := &ss0
	ssEnd := ss[len(ss)-1 : len(ss)-1]
	fmt.Printf("ss: %#v\n", ss)
	fmt.Printf("ss0 = ss[0:0]:%#v\n", ss0)
	fmt.Printf("ss[-1:-1]:%#v\n", ssEnd)
	fmt.Printf("*ss: %v\n", ssp)
	fmt.Printf("*ss0: %v\n", ssp0)
	ss0 = append(ss0, "bingo")
	fmt.Print("appended bingo to ss0\n")
	fmt.Printf("ss: %#v\n", ss)
	fmt.Printf("ss[0:0]: %#v\n", ss0)
	fmt.Printf("ssEnd: %#v\n", ssEnd)
	ss0 = append(ssEnd, "dongo")
	fmt.Printf("append dongo to ssEnd")
	fmt.Printf("ss: %#v\n", ss)
	fmt.Printf("ss[0:0]: %#v\n", ss0)
	fmt.Printf("ssEnd: %#v\n", ssEnd)
	fmt.Printf("Wild that that works. Originally thought would be good for deque but no, the ssEnd is not gonna move with the slice")

}

func TestSliceDelete(t *testing.T) {
	ss := []string{"a", "bee", "sea", "doo"}
	t.Log(len(ss), "| ", ss)
	ss = slices.Delete(ss, 3, 4)

	t.Log(len(ss), "| ", ss)
}
func TestPtrEqual(t *testing.T) {
	ss := testingStringSlice()
	pv1 := &ss[3]
	pv2 := getptrOrClosest(ss, 3)
	t.Logf("pv1 == pv2: %t", pv1 == pv2)
	t.Logf("pv1 == &ss[3]: %t", pv1 == &ss[3])
}
func TestJoinString(t *testing.T) {
	outtext := testingStringSlice()
	t.Log("Direct Slice Print")
	t.Log(outtext)
	t.Log("\nSlice Join No-Separator")
	t.Log(strings.Join(outtext, "\n"))

}

func TestErrorReturn(t *testing.T) {
	echek := retCheck("./notarealfile.fileextension")
	t.Logf("Checked done. Got: %v", echek)
	enochek := retAsIs("./notarealfile.fileextension")
	t.Logf("Unchecked done. Got: %v", enochek)
	if echek != enochek {
		t.Logf("is echek==enochek? %t", echek == enochek)
		t.Logf("(same error type but different instances of error)")
	}
	t.Logf("echek = no file? %t\nenochek = no file? %t", os.IsNotExist(echek), os.IsNotExist(enochek))
}
func TestStringIndexing(t *testing.T) {
	str := "~\\what/$&$#^"
	for i := range len(str) {
		t.Logf("%d) %d %s", i, str[i], string(str[i]))
	}
}

func TestPathChanges(t *testing.T) {
	cwd := getcwd(t)
	lop := func(oname, res string) {
		t.Logf("|%s()-> %s", oname, res)
	}
	pdir := filepath.Dir(cwd)
	gpdir := filepath.Dir(pdir)
	t.Logf("cwd = %s", cwd)
	lop("Base", filepath.Base(cwd))
	lop("Dir", pdir)
	lop("Dir x2", gpdir)
	lop("Clean", filepath.Clean(cwd))
	lop("Ext", filepath.Ext(cwd))
	lop("VolumeName", filepath.VolumeName(cwd))
	lop("FromSlash", filepath.FromSlash(cwd))
	dirr, fll := filepath.Split(cwd)
	t.Logf("|Split()-> %s, %s", dirr, fll)
	gfiles, e := filepath.Glob("*.*")
	if e != nil {
		t.Logf("[glob: %e]", e)
	}
	t.Logf("GLOB:\n------")
	for i, fn := range gfiles {
		t.Logf("(%d) %s", i, fn)
	}

}
func TestPathSplit(t *testing.T) {
	cwd := getcwd(t)
	cwdloc, e := filepath.Localize(filepath.FromSlash(filepath.Clean(cwd)))
	if e != nil {
		t.Logf("localize err: %e", e)
	}
	seg := filepath.SplitList(cwdloc)
	for i, s := range seg {
		t.Logf("%d. %s", i, s)
	}
}
