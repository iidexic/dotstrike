package pops

import "testing"

func TestSysPaths(t *testing.T) {
	home, e := SysHomepath()
	if e != nil {
		t.Error(e)
	}
	t.Logf("home: %s", home)
	config, e := SysConfigpath()
	if e != nil {
		t.Error(e)
	}
	t.Logf("config: %s", config)
	cache, e := SysCachepath()
	if e != nil {
		t.Error(e)
	}
	t.Logf("cache: %t", cache)
}
