package pops

import (
	"fmt"
	"os"
)

var (
	HomePath    *string          = nil
	ConfigPath  *string          = nil
	CachePath   *string          = nil
	CWD         *string          = nil
	OsDirErrors map[string]error = make(map[string]error)
)

var (
	ErrorUserDirs  = fmt.Errorf(`Failed to get system config and home directories (not set in env)`)
	ErrorUserDir   = fmt.Errorf(`Failed to get user directory (not set in env)`)
	ErrorHomeDir   = fmt.Errorf(`Failed to get home directory (not set in env)`)
	ErrorConfigDir = fmt.Errorf(`Failed to get config directory (not set in env)`)
)

func SysHomepath() (string, error) {
	return runSysDir(HomePath, "home", os.UserHomeDir)
}

func SysConfigpath() (string, error) {
	return runSysDir(ConfigPath, "config", os.UserConfigDir)
}

func SysCWD() (string, error) {
	return runSysDir(CWD, "cwd", os.Getwd)
}

func SysCachepath() (string, error) {
	return runSysDir(CachePath, "cache", os.UserCacheDir)
}

func runSysDir(pstore *string, name string, f func() (string, error)) (string, error) {
	if pstore != nil && *pstore != "" {
		return *pstore, nil
	}
	str, e := f()
	if e != nil {
		OsDirErrors[name] = e
		return *pstore, ErrorUserDir
	} else if str == "" {
		OsDirErrors[name] = ErrEmptyHome
		return *pstore, ErrEmptySystem
	} else {
		pstore = &str
		return *pstore, nil

	}
}

// ok so:
// 1. Trigger/Check we have sys paths
// 2. make config file path(s)
