package main

import (
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
	//websocket
	e.SetHandle("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.RunServer(w,r)
	})
	e.Run()
}