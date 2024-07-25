package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	lua "github.com/yuin/gopher-lua"
)

// 返回值表示返回值的数量
func head(l *lua.LState) int {
	var (
		host string
		port uint64
		path string
		resp *http.Response
		err  error
		url  string
	)

	// lua的索引从1开始
	host = l.CheckString(1)
	port = uint64(l.CheckInt64(2))
	path = l.CheckString(3)
	url = fmt.Sprintf("http://%s:%d/%s", host, port, path)
	if resp, err = http.Head(url); err != nil {
		// 传递一个满足接口类型的 lua.Value的对象将值推送给lua.LState
		l.Push(lua.LNumber(0))                                     // http状态码
		l.Push(lua.LBool(false))                                   // 确定服务器是否需要基本身份验证
		l.Push(lua.LString(fmt.Sprintf("Request failed:%s", err))) // 错误消息
		return 3
	}

	l.Push(lua.LNumber(resp.StatusCode))
	l.Push(lua.LBool(resp.Header.Get("WWW-Authenticate") != ""))
	l.Push(lua.LString(""))
	return 3
}

func get(l *lua.LState) int {
	var (
		host     string
		port     uint64
		username string
		password string
		path     string
		resp     *http.Response
		err      error
		url      string
		client   *http.Client
		req      *http.Request
	)

	host = l.CheckString(1)
	port = uint64(l.CheckInt64(2))
	username = l.CheckString(3)
	password = l.CheckString(4)
	path = l.CheckString(5)
	url = fmt.Sprintf("http://%s:%d/%s", host, port, path)

	client = new(http.Client)
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		l.Push(lua.LNumber(0))
		l.Push(lua.LBool(false))
		l.Push(lua.LString(fmt.Sprintf("Unable to build GET request: %s", err)))
		return 3
	}

	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}
	if resp, err = client.Do(req); err != nil {
		l.Push(lua.LNumber(0))
		l.Push(lua.LBool(false))
		l.Push(lua.LString(fmt.Sprintf("Unable to send GET request:%s", err)))
		return 3
	}
	l.Push(lua.LNumber(resp.StatusCode))
	l.Push(lua.LBool(false))
	l.Push(lua.LString(""))
	return 3
}

const (
	LuaHttpTypeName = "http"
)

func register(l *lua.LState) {
	// 标识在lua中创建的名称空间
	mt := l.NewTypeMetatable(LuaHttpTypeName)
	l.SetGlobal("http", mt)                     // 全局名称
	l.SetField(mt, "head", l.NewFunction(head)) // 注册静态函数
	l.SetField(mt, "get", l.NewFunction(get))
}

const (
	PluginDir = "../../plugins"
)

func main() {
	var (
		l     *lua.LState
		files []os.FileInfo
		err   error
		f     string
	)

	l = lua.NewState()
	defer l.Close()

	register(l)
	if files, err = ioutil.ReadDir(PluginDir); err != nil {
		log.Fatalln(err)
	}

	for idx := range files {
		fmt.Println("Found plugin: " + files[idx].Name())
		f = fmt.Sprintf("%s/%s", PluginDir, files[idx].Name())
		if err = l.DoFile(f); err != nil {
			log.Fatalln(err)
		}
	}
}
