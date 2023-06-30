package ui

import (
	"strings"
	"time"
)

type UITimePicker struct {
	*ElemBase
	//y代表年M代表月,以此类推,例如: yyyy-MM-dd HH:mm:ss yyyy年M月 yyyy年的M月某天晚上，大概H点 dd/MM/yyyy ...
	Format string
	/*year	年选择器	只提供年列表选择
	  month	年月选择器	只提供年、月选择
	  date	日期选择器	可选择：年、月、日。type默认值，一般可不填
	  time	时间选择器	只提供时、分、秒选择
	  datetime	日期时间选择器	可选择：年、月、日、时、分、秒
	*/
	DisplayType string //默认datetime
	Value       int64  //默认显示时间 1970年以来的ms数
	Text        string
}

var HtmlTimePicker = `<div class="layui-form-item">
<input type="text" class="layui-input" id="{{.Id}}">
<script>
	layui.use(['form','laydate'], function () {
	     var $ = layui.jquery,
	         form = layui.form,
			 laydate = layui.laydate;

	    //日期
        laydate.render({
            elem: '#{{.Id}}'
			,type: '{{.DisplayType}}'
            ,format: '{{.Format}}'
            ,value: new Date({{if gt .Value 0}}{{.Value}}{{end}})
            ,isInitValue: true
            ,done: function(value, date, endDate){
			    //console.log(value); //得到日期生成的值，如：2017-08-18
			}
        });
	});
</script>
</div>`

// y代表年M代表月,以此类推,例如: yyyy-MM-dd HH:mm:ss yyyy年M月 yyyy年的M月某天晚上，大概H点 dd/MM/yyyy ... 1970年以来的ms数
func newTimePicker(e *ElemBase, format, displaytype string, val int64) *UITimePicker {
	if displaytype == "" {
		displaytype = "datetime"
	}
	p := &UITimePicker{e, format, displaytype, val, ""}
	p.self = p
	return p
}
func NewTimePicker(format, displaytype string, val int64) *UITimePicker {
	return newTimePicker(newElem("timepicker", HtmlTimePicker), format, displaytype, val)
}
func (p *UITimePicker) Clone() HtmlElem {
	np := newTimePicker(cloneElem(p.Id, "timepicker", HtmlTimePicker), p.Format, p.DisplayType, p.Value)
	np.ElemBase.clone(p.ElemBase)
	return np
}
func (p *UITimePicker) SetValue(v string) {
	p.Text = v

	f := p.Format
	f = strings.Replace(f, "yyyy", "2006", -1)
	f = strings.Replace(f, "MM", "01", -1)
	f = strings.Replace(f, "dd", "02", -1)
	f = strings.Replace(f, "HH", "15", -1)
	f = strings.Replace(f, "mm", "04", -1)
	f = strings.Replace(f, "ss", "05", -1)
	t, e := time.ParseInLocation(f, v, time.Local)
	if e != nil {
		return
	}
	p.Value = t.Unix() * 1000
}
