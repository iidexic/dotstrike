package uout

import (
	"fmt"
	"reflect"
	"strings"
)

// EZout simplifies constructing strings intended for user output.
// It will automatically add new lines and (manually-set) indentation
//
// Condensed List of Methods & how they extend the string:
//   - Ln(s): if 1+ strings are passed, adds each on newln. if given 0 strings, it adds newln
//   - A(s): adds s, without  newlines or indentation
//   - AF(s, a...): Sprintf's onto String without newln.
//   - V(a): %+v(a) on newln,
//   - NfnV(a): v(a) (no fieldnames) on newln
//   - NV(name, a): '%s: %v'(name, a) on newln (~named var)
//   - F(s, a...): formatter, add Sprintf(s,a...) on newln
//   - ILns(l) adds a numbered list of the  strings passed
//     *
//     *
//     *
//     *

// TODO:(low) Find out if there's a nicer way to handle the issue of indentation - make line, flatten line?

func ezunwrap(a any, index bool, ez *EZout) {
	rv := reflect.ValueOf(a)
	if rk := rv.Kind(); rk == reflect.Slice {
		if index {
			for i := 0; i < rv.Len(); i++ {
				ez.post(fmt.Sprintf("[%02d] %+v", i, rv.Index(i).Interface()))
				//ez.F("[%02d] %+v", i, rv.Index(i).Interface())
			}
		} else {
			for i := 0; i < rv.Len(); i++ {
				ez.post(fmt.Sprintf("%+v", rv.Index(i).Interface()))
			}
		}

	} else if rk == reflect.Map {
		miter := rv.MapRange()
		for miter.Next() {
			ez.F("[%v]: %+v", miter.Key(), miter.Value())
		}
	} else {
		ez.V(a)
	}
}

func flatunwrap(a any, index bool, ez *EZout) {
	rv := reflect.ValueOf(a)
	if rk := rv.Kind(); rk == reflect.Slice {
		ez.A(" (")
		if index {
			for i := 0; i < rv.Len(); i++ {
				ez.AF(" %02d:%v,", i, rv.Index(i).Interface())
			}
		} else {
			for i := 0; i < rv.Len(); i++ {
				ez.AF(" %+v,", rv.Index(i).Interface())
			}
		}
		ez.A(")")
	} else if rk == reflect.Map {
		miter := rv.MapRange()
		ez.A(" {")
		for miter.Next() {
			ez.AF(" %v:%+v,", miter.Key(), miter.Value())
		}
		ez.A("}")
	} else {
		ez.V(a)
	}
}

func flattenLines(ss []string) []string {
	nsl := make([]string, 0, len(ss)*3)
	for _, s := range ss {
		nnewln := strings.Count(s, "\n")
		if nnewln == 0 {
			nsl = append(nsl, s)
			continue
		}
		nsl = append(nsl, strings.Split(s, "\n")...)
	}
	return nsl
}

/* The starting to go crazy section
type ( ezFuncID  int; OutSelect map[string]EZout;
	EZmanager struct { OutSelect; ezbuilder }
	ezbuilder struct { *EZout; ops []ezOp }
	ezOp struct { vcount int; id ezFuncID })
// why am I keeping this
func (o OutSelect) Gmake() *EZmanager { return &EZmanager{OutSelect: o} }
*/

type EZout struct {
	string
	Ind   int
	clear bool
}

func NewOut(s string) EZout {
	return EZout{string: s, Ind: 0}
}

func NewOutf(s string, a ...any) EZout {

	return EZout{string: fmt.Sprintf(s, a...), Ind: 0}
}
func (E *EZout) WipeOnOutput(b bool) *EZout {
	E.clear = b
	return E
}

func (E EZout) String() string {
	if E.clear {
		defer func() {
			E.string = ""
		}()
	}
	return E.string

}

// pre adds newline and indentation
func (E *EZout) pre() {
	E.string += E.getpre()
}

func (E *EZout) getpre() string {
	out := "\n"
	if E.Ind > 0 {
		return out + E.getInd()
	}
	return out
}

func (E *EZout) getInd() string { return indents(E.Ind) }

// post iterates over lines in s and applies them to the string with proper indentation
func (E *EZout) post(s string) {
	strings.SplitSeq(s, "\n")(func(ln string) bool {
		E.pre()
		E.string += ln
		return true
	})
}

func indents(n int) string { return strings.Repeat("	", n) }

// Ln adds one or more strings, each on a new line
//
// If no strings are passed, it adds a new line. Indents will not be added.
func (E *EZout) Ln(s ...string) {
	if len(s) == 0 {
		E.string += "\n"
	}
	for _, txt := range s {
		E.pre()
		E.string += txt
	}
}

func (E *EZout) PrefixF(ind int, s string, a ...any) {
	prefix := ""
	if ind > 0 {
		prefix = E.getInd()
	}
	E.string = prefix + fmt.Sprintf(s, a...) + E.string
}

func (E *EZout) PrefixV(ind int, a any) {
	prefix := ""
	if ind > 0 {
		prefix = E.getInd()
	}
	E.string = prefix + fmt.Sprintf("%+v", a) + E.string
}

func (E *EZout) LnSplit(s string) {
	lns := strings.Split(s, "\n")
	switch {
	case len(lns) == 1 && lns[0] != "":
		E.V(lns[0])
	case len(lns) > 1:
		E.Ln(lns...)
	}
}

// A directly adds s, without space/newline
func (E *EZout) A(s string) {
	E.string += s
}

// F formats s on a new line
func (E *EZout) F(s string, a ...any) {
	E.pre()
	E.string += fmt.Sprintf(s, a...)
}

func (E *EZout) AF(s string, a ...any) {
	E.string += " " + fmt.Sprintf(s, a...)
}

// ILns (Indexed Lines) prints a numbered list of strings in l, each on a new line
func (E *EZout) ILns(l []string) {
	l = flattenLines(l)
	for i, s := range l {
		E.F("[%02d] %s", i, s)
	}
}

// V adds 1 value with %+v on a new line
func (E *EZout) V(a any) {
	E.pre()
	E.string += fmt.Sprintf("%+v", a)
}

// NfnV (No field name) prints with %v on a new line
func (E *EZout) NfnV(a any) {
	E.pre()
	E.string += fmt.Sprintf("%v", a)
}

// NV adds a named val ("name: a") on a new line
func (E *EZout) NV(name string, a any) {
	E.pre()
	E.string += fmt.Sprintf("%s: %+v", name, a)
}

// IfV prints a if b and aNot if !b, on a new line. Returns b
//
// Always prints a new line
func (E *EZout) IfV(b bool, a, aNot any) bool {
	E.pre()
	if b { // if sa, ok := a.(string); ok && sa != "" && b {
		E.string += fmt.Sprintf("%+v", a)

	} else { // if sna, ok := a.(string); ok && sna != "" && !b
		E.string += fmt.Sprintf("%+v", aNot)
	}
	return b
}

// IfF adds f(s,a...) if b or f(sNot, aNot...) if !b. Returns b
//
// When an empty string is passed for either s/sNot and that string
// would be added, IfF adds no newline and no text.
//   - i.e. b and s=="" or !b and sNot=="" doesn't change EZ.String
func (E *EZout) IfF(b bool, s, sNot string, a, aNot any) bool {
	if b && s != "" {
		E.pre()
		E.string += fmt.Sprintf(s, a)
	} else if !b && sNot != "" {
		E.pre()
		E.string += fmt.Sprintf(sNot, aNot)
	}
	return b
}

// ILV adds an indexed list of values from sa.
//   - If sa is a slice, it adds each '[i] value' on a new line
//   - If sa is a map, it adds each '[key]: val' on a new line (same as LV)
//   - If sa isn't a slice, IV prints the same as V(sa)
func (E *EZout) ILV(sa any) {
	ezunwrap(sa, true, E)
}

// LV adds a list of values from sa.
//   - If sa is a slice, it adds each value on a new line
//   - If sa is a map, it adds a list of '[key]: val', each on a new line
//   - If sa isn't a slice, IV prints the same as V(sa)
func (E *EZout) LV(sa any) {
	ezunwrap(sa, false, E)
}

// Like ILV but Variadic. Does NOT take structs
func (E *EZout) ILVV(sa ...any) {
	if len(sa) == 1 {
		E.F("%+v", sa[0])
	} else {
		for i, a := range sa {
			E.F("[%d] %+v", i, a)
		}
	}
}

// Prints a list, flattened into a single-line comma-separated list of values, in parentheses
func (E *EZout) FlatLV(sa any) {
	flatunwrap(sa, false, E)
}

func (E *EZout) IStringerV(sa ...fmt.Stringer) {
	for i, a := range sa {
		E.F("[%d] %+v", i, a)
	}
}
func (E *EZout) Sep() {
	E.V("------------------------------")
}

// Indent+1
// Returns ptr to itself for chaining
func (E *EZout) IndR() *EZout {
	E.Ind++
	return E
}

// Indent-1
// Returns ptr to itself for chaining
func (E *EZout) IndL() *EZout {
	if E.Ind > 0 {
		E.Ind--
	}
	return E
}

// Indent to 0
// Returns ptr to itself for chaining
func (E *EZout) Ind0() *EZout {
	E.Ind = 0
	return E
}
