package pops

import (
	"fmt"
	"os"
)

var (
	HomePath    *string          = nil
	ConfigPath  *string          = nil
	CachePath   *string          = nil
	OsDirErrors map[string]error = make(map[string]error)
)

var (
	ErrorUserDirs  = fmt.Errorf(`Failed to get system config and home directories (not set in env)`)
	ErrorUserDir   = fmt.Errorf(`Failed to get user directory (not set in env)`)
	ErrorHomeDir   = fmt.Errorf(`Failed to get home directory (not set in env)`)
	ErrorConfigDir = fmt.Errorf(`Failed to get config directory (not set in env)`)
)

func SysHomepath() (string, error) {
	if HomePath != nil && *HomePath != "" {
		return *HomePath, nil
	}
	if h, e := os.UserHomeDir(); e != nil {
		OsDirErrors["home"] = e
		return h, ErrorUserDir
	} else if h == "" {
		OsDirErrors["home"] = ErrEmptyHome
		return h, ErrEmptyHome
	} else {
		HomePath = &h
		return *HomePath, nil
	}
}

func SysConfigpath() (string, error) {
	if ConfigPath != nil && *ConfigPath != "" {
		return *ConfigPath, nil
	}
	if c, e := os.UserConfigDir(); e != nil {
		OsDirErrors["config"] = e
		return c, ErrorUserDir
	} else if c == "" {
		OsDirErrors["config"] = ErrEmptyHome
		return c, ErrEmptySystem
	} else {
		ConfigPath = &c
		return c, nil
	}
}

func SysCachepath() (string, error) {
	if CachePath != nil && *CachePath != "" {
		return *CachePath, nil
	}
	if c, e := os.UserCacheDir(); e != nil {
		OsDirErrors["cache"] = e
		return c, ErrorUserDir
	} else if c == "" {
		OsDirErrors["cache"] = ErrEmptyHome
		return c, ErrEmptySystem
	} else {
		CachePath = &c
		return c, nil
	}
}

// ok so:
// 1. Trigger/Check we have sys paths
// 2. make config file path(s)
