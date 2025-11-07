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
	t.Logf("cache: %s", cache)
	cwd, e := SysCWD()
	if e != nil {
		t.Error(e)
	}
	t.Logf("cwd: %s", cwd)
}
