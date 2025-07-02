package pops

import (
	"fmt"
	"os"
	"path/filepath"
)

func DetailStatPath(p string) string {
	pabs, e := filepath.Abs(p)
	if e != nil {
		panic(e)
	}
	dat, e := os.Stat(pabs)
	if e != nil {
		panic(e)
	}
	return fmt.Sprintf(
		`%s Stat:
Name:%s Size:%d IsDir:%t
fileMode:%v`, p, dat.Name(), dat.Size(), dat.IsDir(), dat.Mode(),
	)

}
