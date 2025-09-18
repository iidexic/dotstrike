package dscore

import (
	"errors"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	toml2 "github.com/pelletier/go-toml/v2"
)

type PelDecErr struct {
	toml2.DecodeError
}

func (P PelDecErr) Error() string {
	return P.DecodeError.Error()
}

func DecodeTomlDataP(r io.Reader) error {
	//CoreConfig()
	dc := toml2.NewDecoder(r)
	e := dc.Decode(&gd.data)
	if e != nil {
		if errors.Is(e, PelDecErr{}) {
			return fmt.Errorf("Ok Error is a DecodeError")
		}
	}
	return e
}

func burntDecode(tomldata string) error {
	md, err := toml.Decode(tomldata, &gd.data)
	if err != nil && errors.Is(err, toml.ParseError{}) {
		return (fmt.Errorf("YES ITS A PARSE ERROR JESUS: %w", err))
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
