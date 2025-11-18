package dscore

import (
	"fmt"

	pops "iidexic.dotstrike/pathops"
	"iidexic.dotstrike/uout"
)

var (
	ErrorFindToml         = fmt.Errorf(`toml config file not found.`)
	ErrorFindMakeToml     = fmt.Errorf(`Failed writing default config to file; User data file not found and could not be made.`)
	ErrorUserDirs         = pops.ErrorUserDirs
	ErrorHomeDir          = pops.ErrorHomeDir
	ErrorConfigDir        = pops.ErrorConfigDir
	ErrorEmptyToml        = fmt.Errorf(`toml config file is empty`)
	ErrorDecodeToml       = fmt.Errorf(`toml decode error`)
	ErrorPartialUndecoded = fmt.Errorf(`some toml keys were not decoded`)
	ErrorAllUndecoded     = fmt.Errorf(`all toml keys were not decoded`)
)

var initlog = uout.NewOut("Initializer Log:")

type initpath struct {
	path   string
	exists bool
	err    error
	read   *pops.ReadResult
}

func (ip *initpath) ReadFile() ([]byte, error) {
	ip.read = pops.ReadFile(ip.path)
	return ip.read.Contents, ip.read.Fail

}
func (ip *initpath) SetError(e error) { ip.err = e }

type namedpaths map[string]*initpath

func (np namedpaths) String() string {
	out := uout.NewOut("[ Named Paths ]")
	for name, p := range np {
		out.F("- %s path:", name)
		out.IndR().V(p.String())
		out.IndL()
	}
	return out.String()
}
func InitString() string { return initlog.String() }

func (ip *initpath) String() string {
	return fmt.Sprintf("path: %s, exists: %t, err: %s", ip.path, ip.exists, ip.err)
}

func (np namedpaths) Add(name string, path string) {
	if pops.PathExistsUsable(path) {
		np[name] = &initpath{path: path, exists: true}
	} else {
		np[name] = &initpath{path: path, exists: false}
	}
}
func (np namedpaths) AddErr(name string, path string, e error) {
	if e != nil {
		np[name] = &initpath{path: path, exists: false, err: e}
	} else {
		np.Add(name, path)
	}
}

func MakeSysConfigPaths(filename string) namedpaths {
	npaths := make(namedpaths, 3)
	// cache path not needed for now
	//cachegood, ecache := pops.SysCachepath()
	if cpath, e := pops.SysConfigpath(); e != nil {
		npaths.AddErr("config", pops.Joinpath(cpath, globalDirConfigRelative, filename), e)
		initlog.F("config path not found: %s", e.Error())
	} else {
		npaths.Add("config", pops.Joinpath(cpath, globalDirConfigRelative, filename))
	}
	if hpath, e := pops.SysHomepath(); e != nil {
		npaths.AddErr("home", pops.Joinpath(hpath, globalDirHomeRelative, filename), e)
		initlog.F("home path not found: %s", e.Error())
	} else {
		npaths.Add("home", pops.Joinpath(hpath, globalDirHomeRelative, filename))
	}
	return npaths
}

// Performs no checks to the path of the toml file beforehand
func loadConfigToml(path string) (*initializer, error) {
	s, e := pops.PathExists(path)
	if e != nil {
		return nil, e
	}
	ip := initpath{path: path, exists: s}
	I := initializer{
		tomlpaths:     make(namedpaths, 1),
		failpaths:     make([]string, 0),
		SysFileErrors: make(map[string]error),
		filename:      pops.Base(path),
		configpath:    path,
	}
	I.tomlpaths[pops.Base(path)] = &ip

	e = I.populateGlobalData()

	return &I, e
}
func loadConfigFromDir(dirpath string) (*initializer, error) {
	return loadConfigToml(pops.Joinpath(dirpath, globalsFilename))
}

type initializer struct {
	tomlpaths     namedpaths
	failpaths     []string
	SysFileErrors map[string]error
	filename      string
	configpath    string
}

func (I *initializer) String() string {
	out := uout.NewOut("[ Initializer ]")
	out.V(I.tomlpaths)
	out.V("Failed Toml Paths:")
	out.Sub().FlatLV(I.failpaths)
	out.V("SysFileErrors:")
	out.IndR().LV(I.SysFileErrors)
	out.IndL().F("filename: %s", I.filename)
	out.F("configpath: %s", I.configpath)
	return out.String()
}

func (I *initializer) populateGlobalData() error {
	if len(I.tomlpaths) == 0 {
		initlog.F("X init tomlpaths: %v", I.tomlpaths)
		return fmt.Errorf("init tomlpaths: %v", I.tomlpaths)
	}
	for _, p := range I.tomlpaths {
		if data, e := p.ReadFile(); e != nil && e != pops.None {
			p.SetError(e)
			initlog.F("X failed to read toml file: %s", p.path)
		} else if len(data) == 0 {
			p.SetError(ErrorEmptyToml)
			initlog.F("X toml file is empty: %s", p.path)
		} else {
			initlog.F("run gd.decodeAsConfig on %s", p.path)
			e = gd.decodeAsConfig(data)
			if e != nil {
				initlog.F("X error decoding toml file: %s", p.path)
				p.SetError(e)
			}
		}

		if gd.status == success {
			GlobalConfigPath = p.path
			gd.dsconfigPath = p.path
			gd.loaded = true
			return nil
		}
		initlog.F("X failed to decode toml file: %s", p.path)
	}
	return ErrorDecodeToml
}

func (I *initializer) Config() error {
	if I.filename == "" {
		I.filename = globalsFilename
	}
	I.tomlpaths = MakeSysConfigPaths(I.filename)
	if len(I.tomlpaths) == 0 {
		return pops.ErrorUserDirs
	}
	e := I.populateGlobalData()
	if e != nil {

		return e
	}
	return nil
}

func (I *initializer) WriteConfigDefaults() error {
	if len(I.tomlpaths) == 0 {
		return ErrorUserDirs
	}
	var ew error
	for name, p := range I.tomlpaths {
		e := gd.encodeDefaultsTo(p.path)
		if e != nil {
			extenderror(ew, e,
				fmt.Sprintf("write to '%s' path (%s)\n	(previously existed: %t, now exists: %t)",
					name, p.path, p.exists, pops.PathExistsUsable(p.path)))
		} else {
			GlobalConfigPath = p.path
			return nil
		}
	}
	return ew
}

func extenderror(em error, e error, msg string) error {
	switch {
	case em == nil && e != nil:
		return e
	case em != nil && e == nil:
		return em
	case em != nil && e != nil:
		return fmt.Errorf("%w\n%s: %w", em, msg, e)
	}
	return nil

}
