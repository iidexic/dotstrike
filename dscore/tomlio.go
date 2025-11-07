package dscore

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	toml2 "github.com/pelletier/go-toml/v2"
)

func DecodeTomlDataP(r io.Reader) error {
	//CoreConfig()
	dc := toml2.NewDecoder(r)
	e := dc.Decode(&gd.data)
	if e != nil {
		if _, ok := any(e).(toml2.DecodeError); ok {
			return fmt.Errorf("Error is a DecodeError: %w", e)
		}
	}
	return e
}

func burntDecode(tomldata string) error {
	md, err := toml.Decode(tomldata, &gd.data)
	if err != nil {
		if _, ok := any(err).(toml.ParseError); ok {
			return fmt.Errorf("Error is a ParseError: %w", err)
		}
	}
	gd.md = md
	return err
}

func trueBurntDecode(r io.Reader) error {
	dc := toml.NewDecoder(r)
	mdat, e := dc.Decode(&gd.data)
	gd.md = mdat
	return e
}
