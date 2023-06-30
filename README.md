# go-layui

# example
```
package main

import (
	"fmt"

	"github.com/you285168/golayui/app"
	"github.com/you285168/golayui/app/ui"
)

var a *app.App

func main() {
	a = app.NewAdminApp(Config.String("addr", ""), "example", "123")
	a.AddPageGroup(DisPlay())
	a.Run()
}

func DisPlay() *ui.PageGroup {
	gp := ui.NewPageGroup("日志查询", "layui-icon layui-icon-template-1")

	gp.AddFrame(newFrame("/loggrep1", "内网日志", "layui-icon layui-icon-face-cry",
		[]string{"grep?s=SL_SERVER_ID&t=SL_LOG_TYPE&d=SL_TIME_SEL&h=SL_TIME_HOUR&u=SL_ROLE_ID&c=SL_CONTENT>>"}))

	gp.AddFrame(newFrame("/loggrep2", "阿里云日志", "layui-icon layui-icon-face-cry",
		[]string{"grep?s=SL_SERVER_ID&t=SL_LOG_TYPE&d=SL_TIME_SEL&h=SL_TIME_HOUR&u=SL_ROLE_ID&c=SL_CONTENT>>"}))

	gp.AddFrame(newFrame("/loggrep3", "版署", "layui-icon layui-icon-face-cry",
		[]string{"grep?s=SL_SERVER_ID&t=SL_LOG_TYPE&d=SL_TIME_SEL&h=SL_TIME_HOUR&u=SL_ROLE_ID&c=SL_CONTENT>>"}))

	return gp
}

func newFrame(router, tile, icon string, gc []string) *ui.Frame {
	f := ui.NewFrame(router, tile, icon, nil)
	hours := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}

	zones := []string{"1001", "1002"}
	logtypes := []string{"上报日志", "Game日志", "消息日志", "战斗日志", "匹配日志", "流水日志"}
	var uiTimePicker *ui.UITimePicker
	var uiHourSelect *ui.UISelect
	var uiZoneSelect *ui.UISelect
	var uiRadio *ui.UIRadio
	var uiRoleEditer *ui.UIEditer
	var uiContentEditer *ui.UIEditer
	var uiResultText *ui.UIText
	var uiButton *ui.UIButton
	f.AddRow().AddLabel("日期").AddTimePicker("yyyyMMdd", "date", 0, &uiTimePicker).AddLabel("小时").AddSelect(0, hours, &uiHourSelect)
	f.AddRow().AddLabel("服务器列表").AddSelect(0, zones, &uiZoneSelect)
	f.AddRow().AddRadio(0, logtypes, &uiRadio)
	f.AddRow().AddLabel("RoleId").AddEditer("RoleId", "", false, &uiRoleEditer).AddLabel("查询内容").AddEditer("Content", "", false, &uiContentEditer).AddButton("查询", func(username string, param map[string]string, getElem ui.GetElem) *ui.Response {
		fmt.Println(param)

		//uiResultText.SetText("hello world")
		getElem(uiResultText).(*ui.UIText).Text = "hello world"
		return ui.ResponseURL(router)
	}, &uiButton)
	f.AddText("", &uiResultText)
	return f
}
```