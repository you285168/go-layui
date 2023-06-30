package ui

import (
	"strings"
)

type UISpace struct {
	*ElemBase
}

var HtmlSpace = `<div class="layui-form-item">
<label class="layui-form-label">&nbsp;</label>
</div>`

func NewSpace() *UISpace {
	return &UISpace{&ElemBase{}}
}
func (l *UISpace) Type() string {
	return "space"
}
func (l *UISpace) ID() string {
	return ""
}
func (l *UISpace) Clone() HtmlElem {
	nl := NewSpace()
	nl.ElemBase.clone(l.ElemBase)
	return nl
}
func (l *UISpace) Render() string {
	var buf strings.Builder
	buf.WriteString(HtmlSpace)
	return buf.String()
}
