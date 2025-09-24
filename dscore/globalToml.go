package dscore

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"iidexic.dotstrike/config"
	pops "iidexic.dotstrike/pathops"
)

var (
	ErrNotModified    error = fmt.Errorf("Attempted write of un-modified temp data")
	ErrModifiedNoInit       = fmt.Errorf("Attempted write of modified UN-INITIALIZED temp data")
	ErrNoInit               = fmt.Errorf("Attempted write of un-initialized temp data")
)

// For some reason BurntSushi/toml always puts "varname =", even if I'm doing all marshaling
func (p prefs) MarshalTOML() ([]byte, error) {
	if len(p.Bools) == 0 {
		return []byte("{}\n"), nil
	}
	output := make([]string, len(p.Bools))
	i := 0
	for k, v := range p.Bools {
		keyOut := fmt.Sprintf("%s = %t", k.String(), v)
		output[i] = keyOut
		i++
	}
	outstring := fmt.Sprintf("{ %s }", strings.Join(output, ", "))
	return []byte(outstring), nil
}

func (p *prefs) UnmarshalTOML(data any) error {
	anymap, good := data.(map[string]any)
	if !good {
		return fmt.Errorf("Failed to assert data as map[string]any")
	}
	p.Bools = make(map[ConfigOption]bool, len(anymap))

	for k, val := range anymap {
		switch val := val.(type) {
		case bool:
			p.Bools[config.OptFrom(k)] = val
		default:
			return fmt.Errorf("Not good type bad type :(")

		}
	}
	return nil
}

// should only be used when very first writing a non-existent dotstrikeData.toml
func (G *globals) encodeDefaults() error {
	file, e := pops.MakeOpenFileF(globalsFilepath())
	if e != nil {
		return e
	}
	defer file.Close()
	encode := toml.NewEncoder(file)
	e = encode.Encode(G.data)
	if e != nil {
		return e
	} else {
		return nil
	}
}

// encodeModified gm data exclusively to main toml
func (gm *globalModify) encodeModified() error {
	file, e := pops.OpenFileRW(globalsFilepath())
	if e != nil || file == nil {
		return e
	}
	defer file.Close()
	file.Truncate(0)
	encode := toml.NewEncoder(file)
	e = encode.Encode(gm.globalData)
	if e != nil {
		return e
	} else {
		return nil
	}
}
