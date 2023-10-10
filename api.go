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

func HandleMergelay(a *App) {
	a.handler.HandleFunc("/api/mergely", func(w http.ResponseWriter, r *http.Request) {
		params := a.ParseHttpParams(r)
		user := params["username"]
		event_id := params["event_id"]
		url_router := params["url_router"]
		file := params["file"]

		res := &ui.ApiRsp{}
		defer func() {
			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
		}()

		h := a.GetElem(user, url_router, event_id)
		if h == nil {
			res.Code = 1
			res.Msg = "no mergely item id=" + event_id
			return
		}
		mer := h.(*ui.UIMergely)

		if mer.F == nil {
			res.Code = 1
			res.Msg = "you should overload UIMergely.OnGetFile"
		} else {
			res.Msg = mer.F(user, file)
		}
	})
}

func HandleUpload(a *App) {
	a.handler.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		res := &ui.Response{}
		defer func() {
			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
		}()

		params := a.ParseHttpParams(r)
		user := params["username"]
		event_id := params["event_id"]
		url_router := params["url_router"]
		res.RedirectUrl = url_router
		up := a.GetElem(user, url_router, event_id)
		if up == nil {
			res.Error = "no upload item"
			return
		}
		u1 := up.(*ui.UIUpload)
		var fileName string
		var fileData []byte

		reader, err := r.MultipartReader()
		if err != nil {
			res.Error = err.Error()
			return
		}

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				res.Error = err.Error()
				return
			}
			if part.FileName() == "" { // this is FormData
				//data, _ := ioutil.ReadAll(part)
				//fmt.Printf("FormData=[%s]\n", string(data))
			} else { // This is FileData
				fileName = part.FileName()
				data, _ := io.ReadAll(part)
				fileData = append(fileData, data...)
			}
		}

		if u1.OnUpload != nil {
			u1.OnUpload(user, params, fileName, fileData)
		}
	})
}

type TableResponse struct {
	Code  int64         `json:"code"`
	Count int64         `json:"count"`
	data  []interface{} `json:"data"`
	Msg   string        `json:"msg"`
}

func TableAddCols(param map[string]string) []string {
	cols := make([]string, 0)
	for i := 0; i < 1000; i++ {
		is := strconv.Itoa(i)
		if n, ok := param[is]; ok {
			cols = append(cols, n)
		} else {
			return cols
		}
	}
	return cols
}

func HandleTable(a *App) {
	a.handler.HandleFunc("/api/table", func(w http.ResponseWriter, r *http.Request) {
		params := a.ParseHttpParams(r)
		user := params["username"]
		event_id := params["event_id"]
		url_router := params["url_router"]
		tbnew := a.GetElem(user, url_router, event_id)
		if tbnew == nil {
			fmt.Println("tbnew is nil")
			io.WriteString(w, `{"code":1,"msg":"no table","count":0,"data":[]}`)
			return
		}
		tb := tbnew.(*ui.UITable)
		oper := params["oper"]
		if oper == "data" {
			page := params["page"]
			limit := params["limit"]
			p, _ := strconv.Atoi(page)
			l, _ := strconv.Atoi(limit)
			sr := params["key[search]"]
			var cnt int
			var data [][]string
			if tb.FuncData != nil {
				cnt, data = tb.FuncData(user, p, l, sr)
			} else {
				cnt, data = tb.TableGetData(user, p, l, sr)
			}
			retData := ""
			ln := 0
			for _, cols := range data {
				if len(tb.Header) == len(cols) {
					if ln > 0 {
						retData += ","
					}
					cc := "{"
					for i, c := range cols {
						if i > 0 {
							cc += ","
						}
						k := "col" + strconv.Itoa(i)
						kk, _ := json.Marshal(&k)
						v, _ := json.Marshal(&c)
						cc += string(kk) + ":" + string(v)
					}
					cc += "}"
					retData += cc
					ln++
				}
			}
			io.WriteString(w, fmt.Sprintf("{\"code\":0,\"msg\":\"\",\"count\":%d,\"data\":[%s]}", cnt, retData))
		} else if oper == "edit" {
			var res ui.ApiRsp
			if tb.FuncEvent != nil {
				res = tb.FuncEvent(user, ui.TOEdit, TableAddCols(params))
			} else {
				res = tb.TableEvent(user, ui.TOEdit, TableAddCols(params))
			}

			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
		} else if oper == "del" {
			var res ui.ApiRsp
			if tb.FuncEvent != nil {
				res = tb.FuncEvent(user, ui.TODel, TableAddCols(params))
			} else {
				res = tb.TableEvent(user, ui.TODel, TableAddCols(params))
			}

			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
		} else if oper == "add" {
			var res ui.ApiRsp
			if tb.FuncEvent != nil {
				res = tb.FuncEvent(user, ui.TOAdd, TableAddCols(params))
			} else {
				res = tb.TableEvent(user, ui.TOAdd, TableAddCols(params))
			}

			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
		} else if oper == "url" {
			if tb.FuncUrl != nil {
				href := params["href"]
				io.WriteString(w, tb.FuncUrl(user, href))
			} else {
				io.WriteString(w, Html404)
			}
		} else if strings.HasPrefix(oper, "user_") {
			res := &ui.Response{Error: "error param"}
			ev, e := strconv.Atoi(oper[5:])
			if e == nil && tb.FuncUsrEvent != nil {
				res = tb.FuncUsrEvent(user, ev, TableAddCols(params))
			}

			ret_json, _ := json.Marshal(res)
			io.WriteString(w, string(ret_json))
		} else {
			io.WriteString(w, `{"code":1,"msg":"error params","count":0,"data":[]}`)
		}
	})
}
