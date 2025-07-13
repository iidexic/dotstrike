package dscore

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func DumpGlobals() []string {
	dump := []string{
		"__GLOBALS__",
		gd.status.Error(),
		fmt.Sprintf("globals loaded: %t", gd.loaded),
		fmt.Sprintf("preferences: %+v", gd.data.Prefs),
		fmt.Sprintf("globals file path: %s", gd.dsconfigPath),
		fmt.Sprintf("checked paths: %v", gd.checkedpaths),
		"-- user cfgs --",
	}
	for i, c := range gd.data.Cfgs {
		dump = append(dump, fmt.Sprintf("[c%d] %s", i, c.status()))
	}
	dump = append(dump, "__MESSAGES__\n")
	dump = append(dump, gd.GlobalMessage...)
	return dump
}

func (G globals) DumpRaw() string {
	return fmt.Sprintf("%+v", G)
}

func printTkeys(keys []toml.Key) {
	for i, k := range keys {
		fmt.Printf("[%d] %s (%+v)\n", i, k.String(), k)
	}
}
func printcfgs(ptrcfg []cfg) {
	for i, cf := range ptrcfg {
		fmt.Printf("[%d] %+v\n", i, cf)
	}
}

func CheckDataDecode(decoded globalData, md toml.MetaData) {

	keys := md.Keys()
	und := md.Undecoded()
	_ = md
	fmt.Println("Decode Results:")
	print(`╭───────────────╮
│ Metadata keys │
╰───────────────╯
`)
	printTkeys(keys)
	print(`╭─────────────╮	
│  Undecoded  │
╰─────────────╯
`)
	printTkeys(und)
	fmt.Printf(`╭────────╮	
│  Data  │
╰────────╯
targetPath:%s
cfgs:`, decoded.TargetPath)
	printcfgs(decoded.Cfgs)
}
