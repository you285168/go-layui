package app

import . "github.com/you285168/golayui/ui"

func UserSetting(a *App) {
	if !a.admin {
		return
	}
	var uiEditer *UIEditer
	var uiTextArea *UITextArea
	frameUserSetting := NewFrame("/usersetting", "UserSetting", "layui-icon layui-icon-username", func(user string) {
		u := UserMgr.GetUser(user)
		a.GetElem(user, "/usersetting", uiEditer.Id).(*UIEditer).Text = u.Phone
		a.GetElem(user, "/usersetting", uiTextArea.Id).(*UITextArea).Text = u.Remark
	})
	frameUserSetting.AddRow().AddLabel("手机号").AddEditer("", "", false, &uiEditer)
	frameUserSetting.AddLabel("备注").AddTextArea("", "", &uiTextArea)
	frameUserSetting.AddButton("保存", func(user string, params map[string]string, getElem GetElem) *Response {
		u := UserMgr.GetUser(user)
		if u.Name != "" {
			u.Phone = params[uiEditer.Id]
			u.Remark = params[uiTextArea.Id]
			UserMgr.AddUser(u)
			return ResponseError("保存成功")
		}
		return ResponseError("保存失败")
	}, nil)
	a.AddFrame(frameUserSetting)
}
