package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/you285168/go-layui/ui"
)

type UserSign struct {
	Name string `json:"username"`
	Pwd  string `json:"password"`
}

func HandleLogin(a *App) {
	a.handler.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		res := &ui.ApiRsp{}
		ret := 0
		cv := ""
		for {
			if r.Body == nil {
				ret = 2
				break
			}
			u := UserSign{}
			e := json.NewDecoder(r.Body).Decode(&u)
			if e != nil {
				res.Msg = e.Error()
				ret = 3
				break
			}
			user := UserMgr.GetUser(u.Name)
			if user.Name == "" {
				ret = 1
				break
			}
			if u.Pwd != user.Pwd {
				ret = 5
				break
			}
			cv, e = Encrypt(user.Name+"|"+strconv.FormatInt(time.Now().Unix(), 10), UserMgr.key)
			if e != nil {
				res.Msg = e.Error()
				ret = 6
				break
			}
			break
		}

		if ret != 0 {
			res.Code = ret
			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
			return
		}

		cookie := &http.Cookie{
			Name:    "token",
			Value:   cv,
			Expires: time.Now().Add(TOKEN_EXPIRE_TIME),
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		ret_json, _ := json.Marshal(res)
		io.WriteString(w, string(ret_json))
	})
}

func HandleSignUp(a *App) {
	a.handler.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		res := &ui.ApiRsp{}
		ret := 0
		cv := ""
		u := UserSign{}
		for {
			if r.Body == nil {
				ret = 2
				break
			}
			e := json.NewDecoder(r.Body).Decode(&u)
			if e != nil {
				res.Msg = e.Error()
				ret = 3
				break
			}
			if u.Name == "" || u.Pwd == "" {
				ret = 7
				break
			}
			if !a.openRegiste { //admin可以注册
				ret = 8
				res.Msg = "registe close"
				break
			}
			user := UserMgr.GetUser(u.Name)
			if user.Name != "" {
				ret = 1
				break
			}
			cv, e = Encrypt(u.Name+"|"+strconv.FormatInt(time.Now().Unix(), 10), UserMgr.key)
			if e != nil {
				res.Msg = e.Error()
				ret = 6
				break
			}
			break
		}

		if ret != 0 {
			res.Code = ret
			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
			return
		}

		UserMgr.AddUser(User{Name: u.Name, Pwd: u.Pwd})

		cookie := &http.Cookie{
			Name:    "token",
			Value:   cv,
			Expires: time.Now().Add(TOKEN_EXPIRE_TIME),
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		ret_json, _ := json.Marshal(res)
		io.WriteString(w, string(ret_json))
	})
}

type UserChpwd struct {
	OldP string `json:"old_password"`
	NewP string `json:"new_password"`
	AgaP string `json:"again_password"`
}

func HandleChpwd(a *App) {
	a.handler.HandleFunc("/api/chpwd", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		params := a.ParseHttpParams(r)
		name := params["username"]
		user := UserMgr.GetUser(name)
		if name == "" || user.Name == "" {
			w.Write([]byte(LoginHtml))
			return
		}

		res := &ui.ApiRsp{}
		ret := 0
		u := UserChpwd{}
		for {
			if r.Body == nil {
				ret = 1
				break
			}
			json.Unmarshal(body, &u)
			if u.OldP != user.Pwd {
				ret = 2
				break
			}
			if u.NewP == "" {
				ret = 3
				break
			}
			break
		}

		if ret != 0 {
			res.Code = ret
			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
			return
		}

		user.Pwd = u.NewP
		UserMgr.AddUser(user)

		ret_json, _ := json.Marshal(res)
		io.WriteString(w, string(ret_json))
	})
}

func (a *App) restoreElem(params map[string]string) {
	url_router := params["url_router"]
	user := params["username"]
	for k, v := range params {
		if strings.HasPrefix(k, ui.HTML_ELEM_ID_HEAD) {
			elem := a.GetElem(user, url_router, k)
			switch e := elem.(type) {
			case *ui.UIText:
				e.Text = v
			case *ui.UIEditer:
				e.Text = v
			case *ui.UISelect:
				e.SelIndex, _ = strconv.Atoi(v)
			case *ui.UIRadio:
				e.SelIndex, _ = strconv.Atoi(v)
			case *ui.UITimePicker:
				if e.DisplayType == "date" {
					t1, _ := time.Parse("20060102", v)
					e.Value = t1.Unix() * 1000
				}
			}
		}
	}
}

func HandleButtonClick(a *App) {
	a.handler.HandleFunc("/button_click", func(w http.ResponseWriter, r *http.Request) {
		params := a.ParseHttpParams(r)
		fmt.Println(params)
		event_id, ok := params["event_id"]
		url_router, ok1 := params["url_router"]
		user := params["username"]
		res := &ui.Response{Error: "error param"}

		a.restoreElem(params)
		getElem := func(e ui.HtmlElem) ui.HtmlElem {
			return a.GetElem(user, url_router, e.ID())
		}

		for {
			if !ok || !ok1 {
				break
			}
			e := a.GetElem(user, url_router, event_id)
			if e == nil {
				break
			}
			b := e.(*ui.UIButton)
			if b.Event == nil {
				break
			}
			res = b.Event(user, params, getElem)
			break
		}

		ret_json, _ := json.Marshal(res)
		io.WriteString(w, string(ret_json))
	})
}

func addMenu(a *App, p *ui.PageGroup, m *ui.MenuChild, user string) {
	m.Href = ""
	m.Icon = p.Icon
	m.Title = p.Title
	m.Target = "_self"

	for _, f := range p.Frames {
		if a.userValidCheck(user, f.Router) {
			m.Child = append(m.Child, ui.MenuChild{nil, f.Router, f.Icon, "_self", f.Title})
		}
	}
	for _, g := range p.Groups {
		nm := ui.MenuChild{}
		addMenu(a, g, &nm, user)
		m.Child = append(m.Child, nm)
	}
}
func HandleMenu(a *App) {
	a.handler.HandleFunc("/api/init", func(w http.ResponseWriter, r *http.Request) {
		params := a.ParseHttpParams(r)
		res := &ui.Menu{}
		user := params["username"]

		if a.home == nil {
			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
			return
		}

		res.HomeInfo.Title = a.title
		res.HomeInfo.Href = a.home.Router
		res.LogoInfo.Title = a.title
		res.LogoInfo.Image = "/uilib/images/logo.png"
		res.LogoInfo.Href = a.home.Router

		for _, p := range a.group {
			nm := ui.MenuChild{}
			addMenu(a, p, &nm, user)
			res.MenuInfo = append(res.MenuInfo, nm)
		}
		ret_json, _ := json.Marshal(res)
		io.WriteString(w, string(ret_json))
	})
}

func HandleClear(a *App) {
	a.handler.HandleFunc("/api/clear", func(w http.ResponseWriter, r *http.Request) {
		params := a.ParseHttpParams(r)
		user := params["username"]
		if user == "" {
			res := &ui.ApiRsp{1, "清理缓存失败"}
			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
			return
		}
		a.Reset(user)
		res := &ui.ApiRsp{0, "清理缓存成功"}
		ret_json, _ := json.Marshal(res)
		io.WriteString(w, string(ret_json))
	})
}
