package ui

type UICheckBox struct {
	*ElemBase
	Text    string
	Checked bool
}

var HtmlCheckBox = `<div class="layui-form-item {{if .Hide}}layui-hide{{end}}">
<input class="layui-form-checkbox" lay-filter="{{.Id}}" {{if .Disable}}disabled=""{{end}} type="checkbox" id="{{.Id}}" title="{{.Text}}" {{if .Checked}}value="1" checked=""{{else}}value="0"{{end}}>
<script>
	layui.use(['form'], function () {
	     var $ = layui.jquery,
	         form = layui.form;
	    
		form.on('checkbox({{.Id}})', function (data) {
		   if(data.elem.checked){
				$("#{{.Id}}").val('1');
			}else{
				$("#{{.Id}}").val('0');
			}
	   	});
	});
</script>
</div>`

func newCheckBox(e *ElemBase, text string, checked bool) *UICheckBox {
	c := &UICheckBox{e, text, checked}
	c.self = c
	return c
}
func NewCheckBox(text string, checked bool) *UICheckBox {
	return newCheckBox(newElem("checkbox", HtmlCheckBox), text, checked)
}
func (c *UICheckBox) Clone() HtmlElem {
	nc := newCheckBox(cloneElem(c.Id, "checkbox", HtmlCheckBox), c.Text, c.Checked)
	nc.ElemBase.clone(c.ElemBase)
	return nc
}
func (c *UICheckBox) SetValue(v string) {
	if v == "1" {
		c.Checked = true
	} else if v == "0" {
		c.Checked = false
	}
}
