package ui

type ApiRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Response struct {
	Error        string //弹出错误
	RedirectUrl  string //执行跳转页面
	ShowInDialog string //跳转页面是否展示在弹出框里 弹出框标题
	ShowHtml     string //直接弹窗显示返回的string
	SelfClose    bool   //是否关闭当前页面
}

func ResponseError(err string) *Response {
	return &Response{Error: err}
}

func ResponseURL(url string) *Response {
	return &Response{RedirectUrl: url, ShowInDialog: "", SelfClose: true}
}

func ResponseDialog(url, title string) *Response {
	return &Response{RedirectUrl: url, ShowInDialog: title}
}

func ResponseHtml(html, title string) *Response {
	return &Response{ShowHtml: `<pre style="word-wrap: break-word; white-space: pre-wrap;">` + html + "</pre>", ShowInDialog: title}
}
