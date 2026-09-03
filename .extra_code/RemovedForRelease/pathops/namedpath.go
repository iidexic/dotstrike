package pops

import "iidexic.dotstrike/uout"

type NamedPath struct {
	Name string
	Path string
}

type NamedPaths map[string]string

func (N NamedPaths) Add(name, path string) {
	N[name] = path
}
func (N NamedPaths) Get(name string) NamedPath {
	return NamedPath{Name: name, Path: N[name]}
}

func (N NamedPaths) String() string {
	if len(N) > 0 {
		if len(N) == 1 {
			for k, v := range N {
				return k + ": " + v
			}
		}
		out := uout.NewOut("NamedPaths")
		out.LV(N)
		return out.String()
	}
	return "NamedPaths: empty"
}
