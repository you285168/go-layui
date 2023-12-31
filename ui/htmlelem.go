package ui

import (
	"bytes"
	"html/template"
	"strconv"
	"sync/atomic"
)

type HtmlElem interface {
	ID() string
	Type() string
	Clone() HtmlElem
	Render() string
	Route() string
	SetRouter(r string)
	SetValue(v string)
}

type ElemBase struct {
	Id      string
	Typ     string
	Html    string
	Align   string
	Rout    string
	Value   string
	Hide    bool
	Disable bool
	self    interface{}
}

var HTML_ELEM_ID uint32 = 0
var HTML_ELEM_ID_HEAD string = "UICtrl"

func newElemWithoutID(typ, html string) *ElemBase {
	return &ElemBase{"", typ, html, "center", "", "", false, false, nil}
}

func newElem(typ, html string) *ElemBase {
	id := atomic.AddUint32(&HTML_ELEM_ID, 1)
	return &ElemBase{HTML_ELEM_ID_HEAD + strconv.Itoa(int(id)), typ, html, "center", "", "", false, false, nil}
}

func cloneElem(id, typ, html string) *ElemBase {
	return &ElemBase{id, typ, html, "center", "", "", false, false, nil}
}

func (e *ElemBase) ID() string {
	return e.Id
}
func (e *ElemBase) Type() string {
	return e.Typ
}
func (e *ElemBase) Route() string {
	return e.Rout
}
func (e *ElemBase) SetRouter(r string) {
	e.Rout = r
}
func (e *ElemBase) SetValue(v string) {
	e.Value = v
}
func (e *ElemBase) clone(oe *ElemBase) {
	e.Align = oe.Align
	e.Hide = oe.Hide
	e.Disable = oe.Disable
}
func (e *ElemBase) Render() string {
	te := template.New("htmlelem")
	te = te.Funcs(template.FuncMap{"RawString": RawString, "RawInt": RawInt, "IntSliceElem": IntSliceElem, "StringSliceElem": StringSliceElem, "Render": Render})
	t, er := te.Parse(e.Html)
	if er != nil {
		return er.Error()
	}
	buf := bytes.NewBufferString("")
	er = t.Execute(buf, e.self)
	if er != nil {
		return er.Error()
	}
	return buf.String()
}

func StringSliceElem(s []string, i int) string {
	if len(s) == 0 {
		return ""
	}
	if i >= len(s) {
		i = 0
	}
	return s[i]
}

func IntSliceElem(s []int, i int) int {
	if len(s) == 0 {
		return 0
	}
	if i >= len(s) {
		i = 0
	}
	return s[i]
}

func RawInt(i int) template.JS {
	return template.JS(strconv.Itoa(i))
}
func RawString(s string) template.JS {
	return template.JS(s)
}
func Render(e HtmlElem) template.HTML {
	return template.HTML(e.Render())
}
