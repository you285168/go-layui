package ui

import (
	"strconv"
)

type UISelect struct {
	*ElemBase
	SelIndex int
	Option   []string
}

var HtmlSelect = `<div class="layui-form-item {{if .Hide}}layui-hide{{end}}">
<select id="{{.Id}}" lay-filter="{{.Id}}">
{{range $i,$v:=.Option}}<option {{if eq $i $.SelIndex}}selected{{end}} {{if $.Disable}}disabled=""{{end}} value="{{$i}}">{{$v}}</option>{{end}}
</select>
<script>
	layui.use(['form'], function () {
	     var $ = layui.jquery,
	         form = layui.form;
	    
		form.on('select({{.Id}})', function (data) {
	       //console.log(this.innerText);
	       //console.log(data.value);
	   	});
	});
</script>
</div>`

func newSelect(e *ElemBase, selindex int, options []string) *UISelect {
	no := make([]string, len(options), len(options))
	copy(no, options)
	s := &UISelect{e, selindex, no}
	s.self = s
	return s
}

func NewSelect(selindex int, options []string) *UISelect {
	return newSelect(newElem("select", HtmlSelect), selindex, options)
}
func (s *UISelect) Clone() HtmlElem {
	ns := newSelect(cloneElem(s.Id, "select", HtmlSelect), s.SelIndex, s.Option)
	ns.ElemBase.clone(s.ElemBase)
	return ns
}
func (s *UISelect) Text() string {
	if s.SelIndex < len(s.Option) {
		return s.Option[s.SelIndex]
	}
	return ""
}
func (s *UISelect) SetValue(v string) {
	i, e := strconv.Atoi(v)
	if e != nil {
		return
	}
	if i < len(s.Option) {
		s.SelIndex = i
	}
}
