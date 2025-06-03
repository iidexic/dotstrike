package dscore

type pathType int

//type pathConfig int

const (
	filePath pathType = iota
	dirPath
	//filePattern //eventually; paths with wildcards/regex
)

type pathInfo struct {
	path  string
	ptype pathType
}

type pathComponent struct {
	pathInfo
	abspath string
	ignores []string
	alias   string
}

type configReadResult int
type configType int

const (
	noInit configReadResult = iota
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
// dsCfg for a given app or location
// stores info on
type dsCfg struct {
	name         string
	sources      []pathComponent
	destinations []pathComponent

}

// from cmd/config.go
type globalConfigLoad struct {
	status configReadResult
	loaded bool
	pathGlobal string
}

// config holds configuration status and data
type config struct {
	status configReadResult
	loaded bool
	cfpath string
	dpaths []string
	data   any
}
func (c *dsCfg) ReadConfigDefaultLocation(){
	for _,p:= range default	
}
/* func (c *Config) ReadCfgCustomLocaton(cfgpath string){

} */


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
