package pops

import (
	"testing"
)

func TestCopyDir(t *testing.T) {
	srcDir := "D:/coding/exampleFiles"
	outDir := "C:/dev/.test_data/file_operations"
	//NOTE: Testing makeRootSubdir Code; normally just write it out. Path doesn't need to exist
	cm := GetCopierMaschine()
	cm.NewJob("test_examplefiles", srcDir, outDir)
	tcopy := cm.GetJob("test_examplefiles")
	err := tcopy.Run()
	if err != nil {
		t.Errorf("COPY ERROR: %v", err)
		t.Logf("[COPYJOB: %+v]", tcopy)
	}
}
