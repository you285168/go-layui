// radio.go
package ui

import (
	"strconv"
)

type UIRadio struct {
	*ElemBase
	SelIndex int //默认选择0
	Option   []string
	Disable  []int
}

var HtmlRadio = `<div class="layui-form-item {{if .Hide}}layui-hide{{end}}">
<div class="layui-inline">{{$ID:=.Id}}{{$SelIdx:=.SelIndex}}
<input type="text" class="layui-hide" id="{{$ID}}" value="{{$SelIdx}}">
{{range $i,$v:=.Option}}<input type="radio" name="{{$ID}}" lay-filter="{{$ID}}" value="{{$i}}" title="{{.}}" {{if eq $i $SelIdx}}checked=""{{end}} 
{{$ty:=IntSliceElem $.Disable $i}} {{if eq $ty 1}} disabled=""{{end}}>{{end}}
</div>
<script>
	layui.use(['form'], function () {
	     var $ = layui.jquery,
	         form = layui.form;
	    //此处即为 radio 的监听事件
	    form.on('radio({{$ID}})', function(data){
	        $("#{{$ID}}").val(data.value)
	    });
	});
</script>
</div>`

func newRadio(e *ElemBase, selindex int, options []string) *UIRadio {
	no := make([]string, len(options), len(options))
	copy(no, options)
	r := &UIRadio{e, selindex, no, make([]int, len(options), len(options))}
	r.self = r
	return r
}

func NewRadio(selindex int, options []string) *UIRadio {
	return newRadio(newElem("radio", HtmlRadio), selindex, options)
}
func (r *UIRadio) Clone() HtmlElem {
	nr := newRadio(cloneElem(r.Id, "radio", HtmlRadio), r.SelIndex, r.Option)
	copy(nr.Disable, r.Disable)
	nr.ElemBase.clone(r.ElemBase)
	return nr
}
func (r *UIRadio) SetDisable(idx []int) {
	for _, i := range idx {
		if i < len(r.Disable) {
			r.Disable[i] = 1
		}
	}
}
func (r *UIRadio) Text() string {
	if r.SelIndex < len(r.Option) {
		return r.Option[r.SelIndex]
	}
	return ""
}
func (r *UIRadio) SetValue(v string) {
	i, e := strconv.Atoi(v)
	if e != nil {
		return
	}
	if i < len(r.Option) {
		r.SelIndex = i
	}
}
