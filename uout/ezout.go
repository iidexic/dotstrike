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

// ── Reflect slice/map unwrapping ────────────────────────────────────
// NOTE: These are basically just methods, might want to change to actual methods or string functions

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

// what was this even supposed to be?
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
	sub   bool
}

func NewOut(s string) EZout {
	return EZout{string: s, Ind: 0}
}

func NewOutf(s string, a ...any) EZout {

	return EZout{string: fmt.Sprintf(s, a...), Ind: 0}
}

// ── Private Utilities ───────────────────────────────────────────────

// indents returns n*tab spaces
func indents(n int) string { return strings.Repeat("	", n) }

// pre adds newline and indentation
func (E *EZout) pre() {
	E.string += E.getpre()
}

// getpre returns a newline and indentation string
func (E *EZout) getpre() string {
	out := "\n"
	if E.Ind > 0 {
		return out + E.getInd()
	}
	return out
}

// getInd returns indentation string
func (E *EZout) getInd() string { return indents(E.Ind) }

// post iterates over lines in s and applies them to the string with proper indentation
func (E *EZout) post(s string) {
	strings.SplitSeq(s, "\n")(func(ln string) bool {
		E.pre()
		E.string += ln
		return true
	})
}

// NOTE: Privated; requires adding Sub to all newline public methods

// Sub will add an indentation for one newline operation.
//
// Run in place of IndR if the indentation will only be used once
func (E *EZout) indSub() *EZout { // Sub(), change back when implemented
	E.sub = true
	E.Ind++
	return E
}

func (E *EZout) endOp() {
	if E.sub {
		E.Ind--
	}
}

// ── Public Util Methods ─────────────────────────────────────────────

func (E *EZout) WipeOnOutput(b bool) *EZout {
	E.clear = b
	return E
}

// String returns E.string. If E.clear is true, it also empties E.string
func (E *EZout) String() string {
	if E.clear {
		defer E.Clear()
	}
	return E.string
}

// Clear empties E.string
func (E *EZout) Clear() { E.string = E.string[:0] }

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

// ── Public String Builder Methods ───────────────────────────────────

// Ln adds one or more strings, each on a new line.
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

// PrefixF adds a formatted prefix line to E.string, with ind*tab spaces before s
func (E *EZout) PrefixF(ind int, s string, a ...any) {
	prefix := ""
	if ind > 0 {
		prefix = E.getInd()
	}
	E.string = prefix + fmt.Sprintf(s, a...) + E.string
}

// PrefixV adds a prefix line to E.string, with ind*tab spaces before a
func (E *EZout) PrefixV(ind int, a any) {
	prefix := ""
	if ind > 0 {
		prefix = E.getInd()
	}
	E.string = prefix + fmt.Sprintf("%+v", a) + E.string
}

// LnSplit splits s on newlines and adds each line to E.string
// This allows to maintain consistent indentation
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

// AF formats s without a new line, and adds a space before s
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

// H adds a header on a new line with %v, indented one less
//
// Effectively a shortcut for: E.IndL().V(a); E.IndR()
func (E *EZout) H(a any) {
	E.Ind--
	E.pre()
	E.string += fmt.Sprintf("%v", a)
	E.Ind++
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

// IfLN effectively prints a pass/fail list of strings. Each string is printed with f(<s/sNot>,snames[i])
//   - If b[i] is true, prints f(s,snames[i])
//   - If b[i] is false, prints f(sNot,snames[i])
//   - if either s or sNot is empty, prints nothing in that respective case
func (E *EZout) IfLN(b []bool, s, sNot string, snames []string) {
	for i, nm := range snames {
		E.IfF(b[i], s, sNot, nm, nm)
	}
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

// Like ILV but Variadic. Doesn't take structs, doesn't unwrap slices or maps
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

// IStringerV prints a list of Stringer values
func (E *EZout) IStringerV(sa ...fmt.Stringer) {
	for i, a := range sa {
		E.F("[%d] %+v", i, a)
	}
}

// Sep adds a separator line of 30 dashes
func (E *EZout) Sep() { E.V("------------------------------") }

// Sep adds a separator line of 30x r
func (E *EZout) Sepr(r rune) { E.V(strings.Repeat(string(r), 30)) }
