package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/you285168/go-layui/ui"
)

var TOKEN_EXPIRE_TIME = 365 * 24 * time.Hour //默认一年

type AppData struct {
	Created time.Time
	Frames  map[string]*ui.Frame
	cache   map[string]map[string]string
}

func NewAppData() *AppData {
	return &AppData{Created: time.Now(), Frames: make(map[string]*ui.Frame, 0)}
}

type UserAuthCheck func(user, router string) bool

type App struct {
	*AppData    //全局数据
	address     string
	title       string
	handler     *http.ServeMux
	data        map[string]*AppData //用户数据
	group       []*ui.PageGroup     //menu
	home        *ui.Frame           //主页
	debug       bool
	admin       bool          //admin模式
	authCheck   UserAuthCheck //验证用户访问权限
	openRegiste bool          //是否开放注册
}

func NewAPP(addr, title string, debug bool) *App {
	return &App{
		AppData: NewAppData(),
		address: addr,
		handler: &http.ServeMux{},
		data:    make(map[string]*AppData),
		title:   title,
		debug:   debug,
	}
}

func NewAdminApp(addr, title, key string) *App {
	defaultKey := "132353s4te7hsfg3"
	if len(key) < 16 {
		key = key + defaultKey
	}
	key = key[:16]
	UserMgr.SetKey(key)
	return &App{
		AppData:     NewAppData(),
		address:     addr,
		title:       title,
		handler:     &http.ServeMux{},
		data:        make(map[string]*AppData, 0),
		admin:       true,
		openRegiste: true,
	}
}

// 是否开放注册
func (a *App) OpenRegiste(r bool) {
	a.openRegiste = r
}

// 权限管理
func (a *App) SetAuthCheck(f UserAuthCheck) {
	a.authCheck = f
}

// 获取http请求参数，从token中提取原始的username和token
func (a *App) ParseHttpParams(r *http.Request) map[string]string {
	r.ParseForm()
	params := make(map[string]string, 0)
	for k, v := range r.Form {
		params[k] = string(v[0])
	}

	c, e := r.Cookie("token")

	if e == nil {
		s, _ := Decrypt(c.Value, UserMgr.key)
		params["token"] = s

		ss := strings.Split(s, "|")
		if len(ss) > 0 {
			name := ss[0]
			usr := UserMgr.GetUser(name)
			if usr.Name == name {
				params["username"] = name
			}
		}
	}

	return params
}

func (a *App) addGroup(p *ui.PageGroup) {
	for _, f := range p.Frames {
		a.AddFrame(f)
	}
	for _, g := range p.Groups {
		a.addGroup(g)
	}
}

func (a *App) Run() error {
	h := a.handler
	for _, p := range a.group {
		a.addGroup(p)
	}
	UserSetting(a)
	//css js
	h.Handle("/uilib/", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, ".css") {
			resp.Header().Set("content-type", "text/css; charset=utf-8")
		}
		http.StripPrefix("/uilib/", http.FileServer(UILib)).ServeHTTP(resp, req)
	}))
	h.Handle("/mergely", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := a.ParseHttpParams(r)
		extend := "&event_id=" + params["event_id"] + "&url_router=" + params["url_router"]
		w.Write([]byte(ui.MergelyPage(params["fl"]+extend, params["fr"]+extend)))
	}))

	for r, _ := range a.Frames {
		tempr := r
		h.HandleFunc(r, func(w http.ResponseWriter, r *http.Request) {
			params := a.ParseHttpParams(r)
			user := params["username"]

			f := a.GetFrame(user, tempr)
			if f == nil {
				w.Write([]byte(Html404))
			} else {
				if f.OnLoad != nil {
					f.OnLoad(user)
				}
				fr := params["f"]
				if fr == "1" {
					w.Write([]byte(f.RenderFrame()))
				} else {
					w.Write([]byte(f.Render()))
				}
			}
		})
	}
	HandleButtonClick(a)
	HandleTable(a)
	HandleMergelay(a)
	HandleUpload(a)

	if a.admin {
		h.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := a.ParseHttpParams(r)
			user := params["username"]
			if user == "" {
				w.Write([]byte(LoginHtml))
				return
			}
			s := fmt.Sprintf(ui.IndexHtml, a.title, user)
			w.Write([]byte(s))
		}))
		h.Handle("/chpwd", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := a.ParseHttpParams(r)
			user := params["username"]
			if user == "" {
				w.Write([]byte(LoginHtml))
				return
			}
			w.Write([]byte(ChpwdHtml))
		}))
		HandleLogin(a)
		HandleSignUp(a)
		HandleChpwd(a)
	} else {
		h.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := fmt.Sprintf(ui.IndexHtml, a.title, "")
			w.Write([]byte(s))
		}))
	}

	HandleMenu(a)
	HandleClear(a)

	return http.ListenAndServe(a.address, h)
}

// 添加页面
func (a *App) AddFrame(f *ui.Frame) *App {
	if a.home == nil {
		a.home = f
	}
	a.Frames[f.Router] = f
	return a
}

func (a *App) GetElem(user, router, id string) ui.HtmlElem {
	f := a.GetFrame(user, router)
	if f != nil {
		if e, o := f.Events[id]; o {
			return e
		}
	}
	return nil
}

func (a *App) getFrame(user, router string) *ui.Frame {
	ha, ok := a.data[user]
	if !ok {
		ha = NewAppData()
		a.data[user] = ha
	}
	if f, ok := ha.Frames[router]; ok {
		return f
	}
	if f, ok := a.Frames[router]; ok {
		nf := f.Clone().(*ui.Frame)
		ha.Frames[router] = nf
		return nf
	}
	return nil
}

func (a *App) GetFrame(user, router string) *ui.Frame {
	f := a.getFrame(user, router)
	if f == nil {
		return f
	}
	if !a.userValidCheck(user, router) {
		return nil
	}
	return f
}

func (a *App) userValidCheck(user, router string) bool {
	if a.authCheck != nil {
		return a.authCheck(user, router)
	}
	return true
}

// 添加组
func (a *App) AddPageGroup(p *ui.PageGroup) *App {
	a.group = append(a.group, p)
	return a
}

// 重置用户缓存数据
func (a *App) Reset(user string) {
	delete(a.data, user)
}
