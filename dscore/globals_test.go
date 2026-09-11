package dscore

import (
	"testing"
)

func TestCoreConfig(t *testing.T) {
	t.Log(gd.Detail())
	LoadGlobals()
	t.Log(gd.Detail())
}
