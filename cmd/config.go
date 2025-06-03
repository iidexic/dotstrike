/* Copyright © 2025 NAME HERE <EMAIL ADDRESS> */
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	pops "iidexic.dotstrike/pathops"
)

const cfgFile string = "dotstrikemainconfig.toml"

type configStatus int
type configType int

const (
	noInit configStatus = iota
	exeBadToml
	exe
	cacheBadToml
	cache
	homeBadToml
	home
)
const (
	core configType = iota
	dirSrc
	dirDest
)

// Main config variable
var cfg = config{}
var extraConfigs = []config{}

type configContainer interface {
	load()
}

type globalconfig struct {
	status configStatus
	apps   []config
	loaded bool
	cfpath string
	dpaths []string
	data   any
}

// config holds configuration status and data
type config struct {
	status configStatus
	loaded bool
	cfpath string
	dpaths []string
	data   any
}

/* func _initCfg() {
	// 1. Check for an existing coreconfig
	for _, p := range cfg.dpaths {
		fname := path.Join(p, cfgFile)
		print("[[Filepath:", fname, "]]")
		cf := pops.ReadFile(fname)
		if cf.Fail == pops.None && len(cfg.cfpath) == 0 {
			cfg.cfpath = p
			cfg.data = cf.Contents
			// just in case 1 is corrupt, continue to check the loop
		} else if cf.Fail == pops.None {
			//store the etra config location
		}

	}
} */

func (gc globalconfig) load() {
	gc.loaded = true
}

func (c *config) getConfig(dirpath string) bool {
	fpath := path.Join(dirpath, cfgFile)
	fread := pops.ReadFile(fpath)
	if fread.Fail == pops.None {
		if !c.loaded {
			c.data = fread.Contents
			c.load()
			c.cfpath = dirpath
			return true
		} else {
			c.dpaths = append(c.dpaths, dirpath)
		}

	}
	return false
}
func (c config) load() {
	print(c.data)
	c.loaded = true
}
func (c config) get(fpath string) {
	_ = c.getConfig(fpath)
}

func coreConfig() {
	exec, ee := os.Executable()
	cachedir, ec := os.UserCacheDir()
	homedir, eh := os.UserHomeDir()
	ce(ee)
	ce(ec)
	ce(eh)
	ccdirs := []string{exec, cachedir, homedir}
	for _, cloc := range ccdirs {
		_ = cfg.getConfig(cloc)
	}
	cexec := cfg.getConfig(exec)
	ccach := cfg.getConfig(cachedir)
	chome := cfg.getConfig(homedir)
	_, _, _ = cexec, ccach, chome
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "prints  config info",
	Long:  `prints loaded configuration data for dotstrike`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().String("directory", "", "dir")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
