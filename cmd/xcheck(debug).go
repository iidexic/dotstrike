/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
	pops "iidexic.dotstrike/pathops"
	"iidexic.dotstrike/uout"
)

var namedPaths = map[string]string{
	"output":          "d:/coding/exampleFiles/OUTPUT",
	"input":           "d:/coding/exampleFiles/INPUT",
	"data-audio":      "d:/coding/exampleFiles/audio",
	"data-images":     "d:/coding/exampleFiles/imagesets",
	"data-filestruct": "d:/coding/exampleFiles/filestruct",
	"data-samplecode": "d:/coding/exampleFiles/sample_code_text",
}

func getpath(path string) string {
	p, ok := namedPaths[path]
	if !ok {
		return path
	}
	return p
}

// BUG: expandNamedPath Doesn't work
func expandNamedPath(path string) string {
	// for k := range namedPaths {
	// 	if strings.HasPrefix(strings.ToLower(path), k) {
	// 		path = namedPaths[k] + path[len(k):]
	// 	}
	// }
	// return path
	// OR:
	path = pops.CleanPath(path)
	if strings.Contains(path, `\`) {
		b, end, _ := strings.Cut(path, `\`)
		if npb := getpath(b); npb != b {
			path = pops.Joinpath(npb, end)
		}
	} else if np := getpath(path); np != path {
		path = pops.Joinpath(np, path)
	}
	return path
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:    "check",
	Hidden: true,
	Short:  "Check some stuff. Generally just for debug",
	Long: `Check command does buncha stuff depending on args/flags.
	Stuff:
	- temp: prints contents of tempdata struct
	- .: prints contents of current directory
	- cwd: prints current working directory
	- from: prints the path of the file that called the current program
	- toml: prints contents of dscore.GlobalConfigPath
	- ls: prints contents of dscore.GlobalConfigPath
	- exists: prints if a path exists
	- dirs: prints paths of home, config, cache, temp, current, and calledfrom
	- everything: prints all of the above

	flagz:
	- wipe: wipes a named path
	- walk: walks a path
	- parray: prints farg array
	- pslice: prints farg slice
	Namedpaths (for wipe/walk):
	- output: (testing) output path
	- input:  (testing) main input path
	- data-audio: (testing) example audio files
	- data-images: (testing) example image files
	- data-filestruct:  (testing) example unpopulated dir structure
	- data-samplecode: (testing) example sample code
	`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(os.Args[0])
		cmd.Println(fmt.Sprintf("Run %s", os.Args[1:]))
		switch {
		case *checkf.ask:
			conf := askConfirmf("Is this question true")
			fmt.Printf("function gave: %t\n", conf)
			confright := askConfirmf("Was that correct")
			if confright {
				print("thats good")
			} else {
				print("oh no")
			}

		case *checkf.walk:
			if len(args) > 0 {
				cmd.Printf("Walking %s\n", args[0])
				wpath := getpath(args[0])
				dir := pops.ReadDir(wpath)
				print(dir.String())
			}
		case *checkf.wipe:
			if len(args) > 0 {
				cmd.Printf("Wiping %s\n", args[0])
				if _, ok := namedPaths[args[0]]; !ok {
					cmd.Println("Can only wipe named paths.")
				} else {
					wpath := getpath(args[0])
					err := checkf.wipeDir(wpath)
					if err != nil {
						cmd.Printf("Error wiping dir: %s\n", err.Error())
					}
				}
			}
		case len(*checkf.pslice) > 0:
			cmd.Printf("len pslice = %d\n", len(*checkf.pslice))
			cmd.Printf("len parray = %d\n", len(*checkf.parray))
			cmd.Printf("len args = %d\n", len(args))
			cmd.Println("pslice args:\n-----------")
			for i, str := range *checkf.pslice {
				cmd.Printf("[%d] %s\n", i, str)
			}
			printArgs(cmd, args)
		case len(*checkf.parray) > 0:
			cmd.Printf("len pslice = %d\n", len(*checkf.pslice))
			cmd.Printf("len parray = %d\n", len(*checkf.parray))
			cmd.Println("parray args:\n-----------")
			for i, str := range *checkf.parray {
				cmd.Printf("[%d] %s\n", i, str)
			}
			cmd.Println("-----------")
			printArgs(cmd, args)
		case len(args) > 0:
			checkf.runDefault(cmd, args)
		default:
		}

	},
}
var checkf = checkCmdFlags{}

func printArgs(cmd *cobra.Command, args []string) {
	cmd.Printf("len args = %d\nargs:\n-----------\n", len(args))
	for i, a := range args {
		cmd.Printf("[%d] %s\n", i, a)
	}
}

var checkargs = []string{"temp", ".", "cwd", "from", "toml", "ls", "exists", "dirs", "everything"}

func (c *checkCmdFlags) runDefault(cmd *cobra.Command, args []string) {
	var runarg string
	var a1 string
	if len(args) > 0 {
		runarg = args[0]
	}
	if len(args) > 1 {
		a1 = strings.ToLower(args[1])
	}
	switch runarg {
	case "temp":
		cmd.Println(os.TempDir())
	case ".":
		cmd.Println(pops.Abs("."))
	case "cwd":
		cmd.Println(pops.Cwd())
	case "from":
		cmd.Println(pops.CalledFrom())
	case "toml":
		cmd.Println(dscore.GlobalConfigPath)
	case "ls", "dir":
		if a1 != "" {
			cmd.Println(pops.PrintDir(a1))
		} else {
			cmd.Println("ls needs a path")
		}
	case "exists":
		if len(args) > 1 {
			exists, e := pops.PathExists(args[1])
			if e != nil {
				cmd.Printf("Error checking exists: %s\n", e.Error())
			} else {
				cmd.Println(exists)
			}
		} else {
			cmd.Println("exists needs a path")
		}
	case "data", "userdata", "user-data", "globals", "cfg":
		if len(a1) > 0 && a1[0:2] == "nv" {
			cmd.Println(dscore.TempData().Detail(false))
		} else {
			cmd.Println(dscore.TempData().Detail(true))
		}
	case "undecoded":
		if a1 == "type" {
			cmd.Println(dscore.UndecodedType())
		} else {
			cmd.Println(dscore.Undecoded())
		}

	case "undecoded-type":
		cmd.Println(dscore.UndecodedType())
	case "md":
		cmd.Println(dscore.MD())
	case "dirs", "paths", "sysdirs":
		// BUG: Getting nil pointer error
		cmd.Printf("Home: %s\n", *pops.HomePath)
		cmd.Printf("Config: %s\n", *pops.ConfigPath)
		cmd.Printf("Cache: %s\n", *pops.CachePath)
		cmd.Printf("Temp: %s\n", os.TempDir())
		cmd.Printf("Current: %s\n", pops.Cwd())
		cmd.Printf("CalledFrom: %s\n", pops.CalledFrom())
	case "everything":
		out := uout.NewOut("----[ Some Stuff ]----")
		out.F("TempDir: %s", os.TempDir())
		abs, e := pops.Abs(".")
		out.IferF("Abs: %s", abs, e)
		out.F("Cwd: %s", pops.Cwd())
		out.F("CalledFrom: %s", pops.CalledFrom())
		out.F("GlobalConfigPath: %s", dscore.GlobalConfigPath)
		out.Sep()
		td := dscore.TempData()
		out.IfNN(td)
		out.Sep()
		cmd.Println(out.String())
	default:
		cmd.Println("check called, arg not recognized")
		cmd.Println("stuff that does something: ")
		for _, a := range checkargs {
			cmd.Println(a)
		}

	}
}

func (c *checkCmdFlags) wipeDir(path string) error {
	if len(path) < 10 {
		return fmt.Errorf("Path too short don't delete that shit")
	}
	nopaths := make([]string, 0)
	if pops.HomePath != nil {
		nopaths = append(nopaths, *pops.HomePath)
	} else {
		return fmt.Errorf("Home path is nil wtf")
	}
	if pops.ConfigPath != nil {
		nopaths = append(nopaths, *pops.ConfigPath)
	} else {
		return fmt.Errorf("Config path is nil wtf")
	}
	if pops.CachePath != nil {
		nopaths = append(nopaths, *pops.CachePath)
	}
	notinpaths := make([]string, 0)
	if runtime.GOOS == "windows" {
		notinpaths = append(notinpaths, "C:\\Windows")
		notinpaths = append(notinpaths, "C:\\Program Files")
		notinpaths = append(notinpaths, "C:\\Program Files (x86)")
		notinpaths = append(notinpaths, "C:\\ProgramData")
	} else {
		// idk what you people (derogatory) need
	}
	if slices.Contains(nopaths, path) {
		return fmt.Errorf("Path is in the list of paths to not delete")
	}
	for _, p := range notinpaths {
		if strings.Contains(path, p) {
			return fmt.Errorf("Path is in the list of paths to not delete")
		}
	}

	return pops.DeleteDirContents(path)
}

type checkCmdFlags struct {
	show, temp, walk, ask, path, wipe *bool
	parray, pslice                    *[]string
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkf.walk = checkCmd.Flags().BoolP("walk", "w", false, "walk dir")
	checkf.temp = checkCmd.Flags().Bool("temp", false, "show contents of temporary storage struct for changes to user data")
	// StringArray would be set multiple times; one arg per flag
	checkf.parray = checkCmd.Flags().StringArray("parray", []string{}, "print farg array")
	checkf.pslice = checkCmd.Flags().StringSlice("pslice", []string{}, "print farg slice")
	checkf.show = checkCmd.Flags().Bool("show", false, "show exe info")
	checkf.ask = checkCmd.Flags().Bool("ask", false, "check askconfirm")
	checkf.wipe = checkCmd.Flags().Bool("wipe", false, "wipe dir contents")
}
