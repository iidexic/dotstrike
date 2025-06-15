package dscore

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	tA := "bananger jeelly"
	tC := make(map[string]string)
	tC["tee"] = "ee"
	tC["add Prefix"] = "one bananger jeelly"
	tC["add Suffix"] = "bananger jeelly sir or madam"
	tC["add Mid"] = "bananger and/or jeelly"
	tC["add All"] = "can I have one bananger and jeelly please"
	tC["Crnch"] = "bnngr jlly"
	tC["runes-=1"] = "a0m0mfdq iddkkx"
	tC["runes+=1"] = "cbobohfs kffmmz"
	tC["all z equal len"] = "zzzzzzzz zzzzzz"
	tC["all z len-1"] = "zzzzzzzz zzzzz"
	tC["all z len+1"] = "zzzzzzzz zzzzzzz"
	tC["no space"] = "banangerjeelly"
	tC["equal"] = "bananger jeelly"

	for k, v := range tC {
		t.Logf("(%s:) %s vs %s: %d", k, tA, v, strings.Compare(tA, v))
	}
	t.Fail()

}

func TestCompareReverse(t *testing.T) {
	tA := "bananger jeelly"
	tC := make(map[string]string)
	tC["tee"] = "ee"
	tC["add Prefix"] = "one bananger jeelly"
	tC["add Suffix"] = "bananger jeelly sir or madam"
	tC["add Mid"] = "bananger and/or jeelly"
	tC["add All"] = "can I have one bananger and jeelly please"
	tC["Crnch"] = "bnngr jlly"
	tC["runes-=1"] = "a0m0mfdq iddkkx"
	tC["runes+=1"] = "cbobohfs kffmmz"
	tC["all z equal len"] = "zzzzzzzz zzzzzz"
	tC["all z len-1"] = "zzzzzzzz zzzzz"
	tC["all z len+1"] = "zzzzzzzz zzzzzzz"
	tC["no space"] = "banangerjeelly"
	tC["equal"] = "bananger jeelly"

	for k, v := range tC {
		t.Logf("(%s:) %s vs %s: %d", k, v, tA, strings.Compare(v, tA))
	}
	t.Fail()

}

func TestMakeSpace(t *testing.T) {
	m := make([]string, 3)
	ms := make([]string, 0, 3)
	t.Logf("len make([]string,3) = %d", len(m))
	t.Logf("len make([]string,0,3) = %d", len(ms))
	// This errors:
	//ms[0] = "this"
	// space is identified but not allocated
	t.Fail()
}
