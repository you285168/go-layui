// user.go
package app

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

type User struct {
	Name   string
	Pwd    string
	Auth   string
	Phone  string
	Remark string
}

var USER_FILE = "app.users"

type userMgr struct {
	f     *os.File
	key   string //秘钥，用于token加密，密码加密，长度16 24 32
	users map[string]User
	mu    sync.Mutex
}

var UserMgr *userMgr

func init() {
	uf, err := os.OpenFile(USER_FILE, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	UserMgr = &userMgr{f: uf, users: make(map[string]User)}
	cnt := 0
	buf := bufio.NewReader(uf)
	n := 0
	for {
		n++
		line, isPrefix, err := buf.ReadLine()
		for isPrefix {
			var next []byte
			next, isPrefix, err = buf.ReadLine()
			line = append(line, next...)
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		if n == 1 && len(line) > 2 { //has bom?
			if line[0] == 239 && line[1] == 187 && line[2] == 191 {
				line = line[3:]
			}
		}
		cnt++
		ss := strings.Split(string(line), "\t")
		if len(ss) == 5 {
			UserMgr.users[ss[0]] = User{ss[0], ss[1], ss[2], ss[3], ss[4]}
		}
	}
	if cnt == len(UserMgr.users) {
		return
	}
	var wbuf strings.Builder
	for _, u := range UserMgr.users {
		fmt.Fprintf(&wbuf, "%s\t%s\t%s\t%s\t%s\r\n", u.Name, u.Pwd, u.Auth, u.Phone, u.Remark)
	}
	UserMgr.f.WriteString(wbuf.String())

}

func (mgr *userMgr) SetKey(key string) {
	mgr.key = key
}

func (mgr *userMgr) GetUser(name string) User {
	name = base64.URLEncoding.EncodeToString([]byte(name))
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if u, ok := mgr.users[name]; ok {
		return mgr.decodeUser(u)
	}
	return User{}
}

func (mgr *userMgr) AddUser(u User) {
	u = mgr.encodeUser(u)
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.users[u.Name] = u
	var buf strings.Builder
	fmt.Fprintf(&buf, "%s\t%s\t%s\t%s\t%s\r\n", u.Name, u.Pwd, u.Auth, u.Phone, u.Remark)
	mgr.f.WriteString(buf.String())
}

func (mgr *userMgr) AllUser() []string {
	mgr.mu.Lock()
	us := make([]string, 0, len(mgr.users))
	for s, _ := range mgr.users {
		name, _ := base64.URLEncoding.DecodeString(s)
		us = append(us, string(name))
	}
	mgr.mu.Unlock()
	sort.Strings(us)
	return us
}

func (mgr *userMgr) encodeUser(ou User) User {
	var nu User
	p, e := Encrypt(ou.Pwd, mgr.key)
	if e != nil {
		return nu
	}
	nu.Pwd = p
	nu.Name = base64.URLEncoding.EncodeToString([]byte(ou.Name))
	nu.Auth = base64.URLEncoding.EncodeToString([]byte(ou.Auth))
	nu.Phone = base64.URLEncoding.EncodeToString([]byte(ou.Phone))
	nu.Remark = base64.URLEncoding.EncodeToString([]byte(ou.Remark))
	return nu
}

func (mgr *userMgr) decodeUser(ou User) User {
	var nu User
	p, e := Decrypt(ou.Pwd, mgr.key)
	if e != nil {
		return nu
	}
	nu.Pwd = p
	name, _ := base64.URLEncoding.DecodeString(ou.Name)
	nu.Name = string(name)
	auth, _ := base64.URLEncoding.DecodeString(ou.Auth)
	nu.Auth = string(auth)
	phone, _ := base64.URLEncoding.DecodeString(ou.Phone)
	nu.Phone = string(phone)
	remark, _ := base64.URLEncoding.DecodeString(ou.Remark)
	nu.Remark = string(remark)
	return nu
}
