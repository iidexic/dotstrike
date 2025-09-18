package dscore

import (
	"bytes"
	"errors"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestDecodeMainPell(t *testing.T) {
	cfgdir := gd.makeCfgPath(globalDirHomeRelative)
	gotcfg := gd.GetConfig(cfgdir)
	if gotcfg {
		raw := []byte(gd.rawContents)
		bread := bytes.NewReader(raw)
		e := DecodeTomlDataP(bread)
		if e != nil {
			t.Errorf("Decode Error :(\n%v", e)
		}
		t.Log(gd.data.DetailSimple())
		bremain := bread.Len()
		bend, err := bread.ReadByte()
		if err != nil {
			t.Logf("Error from reading remainder:\n%s", err.Error())
		}
		t.Logf("Length of remainder:%d\nRemainder data:%s", bremain, string(bend))
	}
}

func TestDecodeMainBurn(t *testing.T) {
	cfgdir := gd.makeCfgPath(globalDirHomeRelative)
	gotcfg := gd.GetConfig(cfgdir)
	if gotcfg {
		e := burntDecode(gd.rawContents)
		if e != nil {
			if errors.Is(e, toml.ParseError{}) {
				t.Log("is a parse error")
			}
			t.Errorf("Decode Error :(\n%s", e.Error())
		}
		t.Log(gd.data.DetailSimple())
		t.Logf("gd.md:\n%+v", gd.md)
		t.Logf("Undecoded:%v", gd.md.Undecoded())
		// bremain := bread.Len()
		// bend, err := bread.ReadByte()
		// if err != nil {
		// 	t.Logf("Error from reading remainder:\n%s", err.Error())
		// }
		// t.Logf("Length of remainder:%d\nRemainder data:%s", bremain, string(bend))

	}

}

func TestDecodeMainTrueBurn(t *testing.T) {
	cfgdir := gd.makeCfgPath(globalDirHomeRelative)
	gotcfg := gd.GetConfig(cfgdir)
	if gotcfg {
		raw := []byte(gd.rawContents)
		bread := bytes.NewReader(raw)
		e := trueBurntDecode(bread)
		if e != nil {
			pe, ok := any(e).(toml.ParseError)
			if ok {
				t.Log("It actually is a parse error")
				t.Logf("Error Position:%+v", pe.Position)
				t.Logf("LastKey: %v", pe.LastKey)
			}
		}
		t.Log(gd.data.DetailSimple())
		bremain := bread.Len()
		bend, err := bread.ReadByte()
		if err != nil {
			t.Logf("Error from reading remainder:\n%s", err.Error())
		}
		t.Logf("Length of remainder:%d\nRemainder data:%s", bremain, string(bend))

	}
}
