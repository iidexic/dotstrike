//go:build mage
// +build mage

package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	pops "iidexic.dotstrike/pathops"
	"iidexic.dotstrike/uout"
)

type target struct {
	os, arch, file string
}

var outputDir = "bin"
var targets = map[string]target{
	"linux-amd64":   {os: "linux", arch: "amd64", file: "ds-linux-amd64"},
	"linux-arm":     {os: "linux", arch: "arm", file: "ds-linux-arm"},
	"linux-arm64":   {os: "linux", arch: "arm64", file: "ds-linux-arm64"},
	"macos-amd64":   {os: "darwin", arch: "amd64", file: "ds-macos-amd64"},
	"macos-arm64":   {os: "darwin", arch: "arm64", file: "ds-macos-arm64"},
	"windows-amd64": {os: "windows", arch: "amd64", file: "ds-win-x64.exe"},
	"windows-arm64": {os: "windows", arch: "arm64", file: "ds-win-arm64.exe"},
	"windows-arm":   {os: "windows", arch: "arm", file: "ds-win-arm.exe"},
}
var acceptedOS = []string{"linux", "darwin", "windows"}
var acceptedArch = []string{"amd64", "arm64", "arm"}

var building string

// Targets - List all available build targets
func Targets() error {
	out := uout.NewOut(" Build Targets by ID:")
	n := 0
	for target, t := range targets {
		n++
		out.F("  %d. '%s' (os %s - arch %s)", n, target, t.os, t.arch)
	}
	println(out.String())
	return nil
}

func setTargetForBuild() error {
	osname := os.Getenv("GOOS")
	arch := os.Getenv("GOARCH")
	if osname == "" {
		osname = runtime.GOOS
	}
	if arch == "" {
		arch = runtime.GOARCH
	}
	if !slices.Contains(acceptedOS, osname) {
		return fmt.Errorf("GOOS '%s' invalid or not supported", osname)
	}
	if !slices.Contains(acceptedArch, arch) {
		return fmt.Errorf("GOARCH '%s' invalid or not supported", arch)
	}
	if osname == "darwin" {
		osname = "macos"
	}
	building = osname + "-" + arch
	return nil
}

func CurrentEnv() error {
	fmt.Printf("GOOS = %s\n", runtime.GOOS)
	fmt.Printf("GOARCH = %s\n", runtime.GOARCH)
	return nil
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Build - A build step that requires additional params, or platform-specific steps
func Build() error {
	mg.Deps(InstallDeps)
	if building == "" {
		err := setTargetForBuild()
		if err != nil {
			return err
		}
	}
	target, ok := targets[building]
	if !ok {
		return fmt.Errorf("target '%s' not found", building)
	}
	filename := pops.Joinpath(outputDir, target.file)

	fmt.Printf("Building %s...\n", filename)
	cmd := exec.Command("go", "build", "-o", filename, ".")

	return cmd.Run()
}

// BuildTargets - Takes a target name/ID (generally 'os-arch') and builds for that target.
//
// See `mage help` or `mage targets` for a list of available targets.
func Buildfor(target string) error {
	t, ok := targets[target]
	if !ok {
		return fmt.Errorf("invalid target: %s", target)
	}
	os.Setenv("GOOS", t.os)
	os.Setenv("GOARCH", t.arch)
	building = target
	return Build()
}
func prepBuildFor(target string) (func() error, error) {
	t, ok := targets[target]
	if !ok {
		return nil, fmt.Errorf("invalid target: %s", target)
	}
	os.Setenv("GOOS", t.os)
	os.Setenv("GOARCH", t.arch)
	building = target
	return Build, nil

}
func buildTargetWithDeps(target string) error {
	t, ok := targets[target]
	if !ok {
		return fmt.Errorf("invalid target: %s", target)
	}
	os.Setenv("GOOS", t.os)
	os.Setenv("GOARCH", t.arch)
	building = target
	mg.Deps(Build)
	return nil
}

// BuildAll - Builds all target os/arch combinations
func BuildAll() error {
	for target := range targets {
		err := Buildfor(target)
		if err != nil {
			return err
		}
	}
	// return nil
	return nil //fmt.Errorf("not implemented")
}
func buildAllWithDeps() error {
	for target := range targets {
		err := buildTargetWithDeps(target)
		if err != nil {
			return err
		}
	}
	return nil
}

// func Install() error {
// 	println("Building + Installing...")
// 	buildos := os.Getenv("GOOS")
// 	buildarch := os.Getenv("GOARCH")
// 	mg.Deps(Build)
// 	if buildos == "windows" && buildarch == "amd64" {
// 		return pops.CopyFile("./bin/ds-win-x64.exe", "c:/dev/bin/ds.exe")
// 	}
// }

// InstallDeps manages your deps, or running package managers.
// TODO: InstallDeps
func InstallDeps() error {
	return nil
}

// Clean - Deletes all ds-* files from bin
func Clean() error {
	bin := os.DirFS(outputDir)
	fmt.Printf("Cleaning  %s...\n", outputDir)
	if bin == nil {
		fmt.Println("bin  removed (no bin)")
	}
	del := 0
	nd := 0
	err := fs.WalkDir(bin, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if bfn := pops.Base(path); strings.HasPrefix(bfn, "ds-") {
			p := pops.Joinpath(outputDir, path)
			p = pops.MakeAbs(p)
			e := os.Remove(p)
			if e != nil {
				return e
			}
			del++
		} else {
			nd++
		}
		return nil
	})
	if del > 0 {
		fmt.Printf("Removed %d files, leaving %d\n", del, nd)
	}
	return err
}
