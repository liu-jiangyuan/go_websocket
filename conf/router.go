package conf

import (
	"github.com/liu-jiangyuan/go_websocket/controller"
	"github.com/liu-jiangyuan/go_websocket/lib"
	"net/http"
)

func InitRoute () {
	lib.Route.Ws = map[string]func(map[string]interface{}) map[string]interface{}{
		//ping
		"PING": func(param map[string]interface{}) map[string]interface{} {
			return param
		},
		//demo
		"Index": controller.Index,
		"send":controller.Send,
	}

	//两种方法都行，自由选择
	//lib.Route.Http["/"] = func(w http.ResponseWriter, r *http.Request){
	//	w.Write([]byte("hello"))
	//}
	//lib.Route.Http["/test"] = func(w http.ResponseWriter, r *http.Request){
	//	controller.Test(w,r)
	//}
	//lib.Route.Http["ws"] = func(w http.ResponseWriter, r *http.Request) {
	//	lib.RunServer(w,r)
	//}

	lib.Route.Http = map[string]func(w http.ResponseWriter, r *http.Request){
		"/":func(w http.ResponseWriter, r *http.Request){
			w.Write([]byte("hello"))
		},
		"/test": func(w http.ResponseWriter, r *http.Request) {
			controller.Test(w,r)
		},
		"/ws" : func(w http.ResponseWriter, r *http.Request) {
			lib.RunServer(w,r)
		},
	}
}