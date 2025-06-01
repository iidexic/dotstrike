package toml

import (
	"fmt"
	"go/types"
)

// attempt 1
type tomlvar interface {
	string | int | bool | float64
}

type tomltype = types.BasicKind

const (
	Bool   tomltype = types.Bool
	Int    tomltype = types.Int
	Float  tomltype = types.Float64
	String tomltype = types.String
)

// attempt 2
type tomlVariable struct {
	name    string
	vartype tomltype
	str     string
	bool    bool
	int     int
	float   float64
	astr    []string
	abool   []bool
	aint    []int
	afloat  []float64
	atable  []tomlVariable
}

type tomlVarTable struct {
	bools   map[string][]bool
	strings map[string][]string
	ints    map[string][]int
	floats  map[string][]float64
}

type tomlConverter []struct {
	raw []string
}

func ParseToml(rawData []string) {
	for i, v := range rawData {
		fmt.Printf("[%d]: %s", i, v)
	}
}
