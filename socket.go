package main

import (
	"github.com/liu-jiangyuan/go_websocket/controller"
	"github.com/liu-jiangyuan/go_websocket/engine"
	"github.com/liu-jiangyuan/go_websocket/lib/ws"
	"log"
	"net/http"
)

func main()  {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	e := engine.InitEngine()
	e.SetHost("0.0.0.0")
	e.SetPort("8089")
	//http
	e.SetHandle("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})
	e.SetHandle("/test", func(writer http.ResponseWriter, request *http.Request) {
		controller.Test(writer,request)
	})
	//websocket
	e.SetHandle("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.RunServer(w,r)
	})
	e.Run()
}
