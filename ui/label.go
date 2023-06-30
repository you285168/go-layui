package ui

type UILabel struct {
	*ElemBase
	Text string
}

var HtmlLabel = `<div class="layui-form-item {{if .Hide}}layui-hide{{end}}">
<label class="layui-form-label" style="text-align:{{.Align}};">{{RawString .Text}}</label>
</div>`

func NewLabel(text string) *UILabel {
	l := &UILabel{newElemWithoutID("label", HtmlLabel), text}
	l.Align = "left"
	l.self = l
	return l
}

func (l *UILabel) Clone() HtmlElem {
	nl := NewLabel(l.Text)
	nl.ElemBase.clone(l.ElemBase)
	return nl
}
