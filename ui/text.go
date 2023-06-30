package ui

import (
	"fmt"
	"strings"
)

type UIText struct {
	*ElemBase
	Text string
}

func (t *UIText) SetText(s string) {
	t.Text = s
}

/*
var HtmlText = `<div class="layui-form-item %s">
<xmp class="layui-code">%s</xmp>
</div>
`
*/
var HtmlText = `<div class="layui-form-item %s">
<pre class="layui-code">%s</pre>
</div>`

func newText(e *ElemBase, text string) *UIText {
	t := &UIText{e, text}
	t.self = t
	return t
}

func NewText(text string) *UIText {
	return newText(newElem("text", HtmlText), text)
}
func (t *UIText) Clone() HtmlElem {
	nt := newText(cloneElem(t.Id, "text", HtmlText), t.Text)
	nt.ElemBase.clone(t.ElemBase)
	return nt
}
func (t *UIText) Render() string {
	var buff strings.Builder
	hide := ""
	if t.Hide {
		hide = "layui-hide"
	}
	fmt.Fprintf(&buff, HtmlText, hide, t.Text)
	return buff.String()
}
