package ui

type GetElem func(HtmlElem) HtmlElem

type OnButtonClick func(username string, param map[string]string, getElem GetElem) *Response

type UIButton struct {
	*ElemBase
	Text  string
	Event OnButtonClick
}

var HtmlButton = `<div class="layui-form-item {{if .Hide}}layui-hide{{end}}" style="text-align:{{.Align}};">
<button type="button" class="layui-btn layui-btn-primary {{if .Disable}}layui-disabled{{end}}" id="{{.Id}}" {{if not .Disable}}onclick="buttonClick('{{.Id}}')"{{end}}>{{RawString .Text}}</button>
</div>`

func newButton(e *ElemBase, text string, event OnButtonClick) *UIButton {
	b := &UIButton{e, text, event}
	b.self = b
	return b
}

func NewButton(text string, event OnButtonClick) *UIButton {
	return newButton(newElem("button", HtmlButton), text, event)
}
func (b *UIButton) Clone() HtmlElem {
	nb := newButton(cloneElem(b.Id, "button", HtmlButton), b.Text, b.Event)
	nb.ElemBase.clone(b.ElemBase)
	return nb
}
