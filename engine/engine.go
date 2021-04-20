package engine

import (
	"github.com/liu-jiangyuan/go_websocket/conf"
	"github.com/liu-jiangyuan/go_websocket/controller"
	"github.com/liu-jiangyuan/go_websocket/lib"
	"log"
	"net/http"
)

type engine struct {
	Route map[string]func(map[string]interface{}) map[string]interface{}
	Port string
	Host string
}

func InitEngine() *engine {
	return &engine{Route: conf.RouteMap}
}

func (e *engine) AddRoute (path string,callBack func(param map[string]interface{})map[string]interface{}) {
	e.Route[path] = callBack
}

func (e *engine) Run() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		controller.Test(writer,request)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		lib.RunServer(w,r)
	})
	log.Printf("server runing on:%+v;\r\n",e.Host+":"+e.Port)
	err := http.ListenAndServe(e.Host+":"+e.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
