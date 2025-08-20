package dscore

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"iidexic.dotstrike/config"
	pops "iidexic.dotstrike/pathops"
)

type prefstring struct {
	Bools map[string]bool
}

var ( //TODO: align these errors a bit more with whatever standard is
	ErrNotModified    error = fmt.Errorf("Attempted write of un-modified temp data")
	ErrModifiedNoInit       = fmt.Errorf("Attempted write of modified UN-INITIALIZED temp data")
	ErrNoInit               = fmt.Errorf("Attempted write of un-initialized temp data")
)

//var errNoTemp error = errors.New("TempData is not initialized or does not exist")

// TODO: replace
func (G *globals) EncodeIfNeeded(tg *globalModify) error {
	if tempData.initialized && tempData.Modified {
		return tg.encodeModified()
	} else if tempData.initialized {
		return ErrNotModified
	} else if tempData.Modified {
		return ErrModifiedNoInit

	}
	return ErrNoInit
}

type liminalData struct {
	Selected         int           `toml:"SelectedSpec"`
	GlobalTargetPath string        `toml:"targetpath"`
	Prefs            liminalPrefs  `toml:"prefs"`
	Specs            []liminalSpec `toml:"specs"`
}

type liminalPrefs struct {
	Bools map[string]bool
}

type liminalSpec struct {
	Alias      string          `toml:"alias"`      // name, unique
	Sources    []pathComponent `toml:"sources"`    // paths marked as origin points
	Targets    []pathComponent `toml:"targets"`    // paths  marked as destination points
	Ignorepat  []string        `toml:"ignores"`    // ignorepat that apply to all sources
	OverrideOn bool            `toml:"overrideOn"` // enable overrides, prevent Overrides being over-written
	Overrides  liminalPrefs    `toml:"overrides"`  // override global prefs
	Ctype      componentType
}

// Leaving this here just to check
// func (p prefs) MarshalTOML() ([]byte, error) {
// 	if len(p.Bools) == 0 {
// 		return []byte("{}\n"), nil
// 	}
// 	stringmap := make(map[string]bool, len(p.Bools))
// 	for k, v := range p.Bools {
// 		stringmap[k.String()] = v
// 	}
// 	lpref := liminalPrefs{Bools: stringmap}
// 	return toml.Marshal(lpref)
// }
// func (p *prefs) UnmarshalTOML(data any) error {
// 	mapt, ok := data.(map[string]any)
// 	if !ok {
// 		return fmt.Errorf("expected table, got %T", data)
// 	}
//
// 	optBool := make(map[ConfigOption]bool, len(mapt))
// 	for k, v := range mapt {
// 		optKey := GetOption(k)
// 		if optKey == NotAnOption {
// 			return fmt.Errorf("Bad Option key - '%s' not an option name", k)
// 		}
// 		bval, ok := v.(bool)
// 		if !ok {
// 			return fmt.Errorf("expected bool val, got %T:(%v)", v, v)
// 		}
// 		optBool[optKey] = bval
// 	}
// 	p.Bools = optBool
// 	return nil
// }

// For some reason BurntSushi/toml always puts "varname =", even if I'm doing all marshaling
func (p prefs) MarshalTOML() ([]byte, error) {
	if len(p.Bools) == 0 {
		return []byte("{}\n"), nil //WARNING: Probably don't work (why?)
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

	for k, v := range anymap {
		switch v.(type) {
		case bool:
			p.Bools[config.OptFrom(k)] = v.(bool)
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
