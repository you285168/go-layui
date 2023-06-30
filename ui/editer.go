package ui

type UIEditer struct {
	*ElemBase
	Prompt string
	Text   string
	Pass   bool
}

var HtmlLineEdit = `<div class="layui-form-item {{if .Hide}}layui-hide{{end}}">
<input {{if .Disable}}disabled=""{{end}} {{if .Pass}}type="password"{{else}}type="text"{{end}} id="{{.Id}}" placeholder="{{.Prompt}}" class="layui-input" value="{{.Text}}">
</div>`

func newEditer(e *ElemBase, prompt, text string, password bool) *UIEditer {
	l := &UIEditer{e, prompt, text, password}
	l.self = l
	return l
}

func NewEditer(prompt, text string, password bool) *UIEditer {
	l := newEditer(newElem("lineedit", HtmlLineEdit), prompt, text, password)
	l.self = l
	return l
}
func (l *UIEditer) Clone() HtmlElem {
	nl := newEditer(cloneElem(l.Id, "lineedit", HtmlLineEdit), l.Prompt, l.Text, l.Pass)
	nl.ElemBase.clone(l.ElemBase)
	return nl
}
func (l *UIEditer) SetValue(v string) {
	l.Text = v
}
