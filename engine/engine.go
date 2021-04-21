package engine

import (
	"github.com/liu-jiangyuan/go_websocket/lib"
	"log"
	"net/http"
)

type Engine struct {
	Route map[string]func(map[string]interface{}) map[string]interface{}
	Port  string
	Host  string
}

func InitEngine () *Engine {
	return &Engine{
		Route: lib.Route.Ws,
		Port: "8080",
		Host: "127.0.0.1",
	}
}
func (e *Engine) SetPort (port string) {
	e.Port = port
}
func (e *Engine) GetPort () string {
	return e.Host
}

func (e *Engine) SetHost (host string) {
	e.Host = host
}
func (e *Engine) GetHost () string {
	return e.Host
}

func (e *Engine) SetHandle (pattern string,action func(writer http.ResponseWriter, request *http.Request)) {
	lib.Route.Http[pattern] = action
}
func (e *Engine) AddRoute (path string,callBack func(param map[string]interface{})map[string]interface{}) {
	lib.Route.Ws[path] = callBack
}

func GetRoute() map[string]func(map[string]interface{}) map[string]interface{} {
	return lib.Route.Ws
}

func (e *Engine) Run() {
	//InitRoute 配置会覆盖默认
	wr := lib.Route.Ws
	for pattern , action := range wr {
		e.AddRoute(pattern,action)
	}
	hr := lib.Route.Http
	for pattern , action := range hr {
		http.HandleFunc(pattern,action)
	}
	log.Printf("engine:%+v",e)
	log.Printf("server runing on:%+v;\r\n",e.Host+":"+e.Port)
	if err := http.ListenAndServe(e.Host+":"+e.Port, nil);err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
