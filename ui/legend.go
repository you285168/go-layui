package ui

type UILegend struct {
	*ElemBase
	Text string
}

var HtmlLegend = `<fieldset class="layui-elem-field layui-field-title {{if .Hide}}layui-hide{{end}}" style="margin-top: 30px;">
  <legend>{{RawString .Text}}</legend>
</fieldset>`

func NewLegend(text string) *UILegend {
	l := &UILegend{newElemWithoutID("legend", HtmlLegend), text}
	l.self = l
	return l
}
func (l *UILegend) Clone() HtmlElem {
	nl := NewLegend(l.Text)
	nl.ElemBase.clone(l.ElemBase)
	return nl
}
