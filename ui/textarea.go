package ui

type UITextArea struct {
	*ElemBase
	Prompt string
	Text   string
	Rows   int
}

var HtmlTextArea = `<div class="layui-form-item {{if .Hide}}layui-hide{{end}}">
<textarea id="{{.Id}}" {{if .Disable}}disabled=""{{end}} placeholder="{{.Prompt}}" class="layui-textarea" rows="{{.Rows}}">{{RawString .Text}}</textarea>
</div>`

func newTextArea(e *ElemBase, prompt, text string) *UITextArea {
	t := &UITextArea{e, prompt, text, 17}
	t.self = t
	return t
}

func NewTextArea(prompt, text string) *UITextArea {
	return newTextArea(newElem("textarea", HtmlTextArea), prompt, text)
}
func (t *UITextArea) Clone() HtmlElem {
	nt := newTextArea(cloneElem(t.Id, "textarea", HtmlTextArea), t.Prompt, t.Text)
	nt.Rows = t.Rows
	nt.ElemBase.clone(t.ElemBase)
	return nt
}
func (t *UITextArea) SetValue(v string) {
	t.Text = v
}
