package ui

import (
	"strings"
)

type UIRow struct {
	*ElemBase
	Elems []HtmlElem
	*uiFactory
}

func NewRow() *UIRow {
	r := &UIRow{&ElemBase{}, make([]HtmlElem, 0), nil}
	r.uiFactory = NewFactory(r)
	return r
}
func (l *UIRow) AddElem(e HtmlElem) UIFactory {
	l.Elems = append(l.Elems, e)
	return l
}

func (l *UIRow) Type() string {
	return "row"
}
func (l *UIRow) ID() string {
	return ""
}
func (l *UIRow) Clone() HtmlElem {
	nl := NewRow()
	for _, r := range l.Elems {
		nl.AddElem(r.Clone())
	}
	nl.ElemBase.clone(l.ElemBase)
	return nl
}
func (l *UIRow) SetRouter(r string) {
	l.ElemBase.SetRouter(r)
	for _, v := range l.Elems {
		v.SetRouter(r)
	}
}
func (l *UIRow) Render() string {
	var buff strings.Builder
	buff.WriteString(`<div class="layui-form-item">`)

	for _, s := range l.Elems {
		buff.WriteString(`<div class="layui-inline">`)
		buff.WriteString(s.Render())
		buff.WriteString(`</div>`)
	}
	buff.WriteString(`</div>`)
	return buff.String()
}
