package main

import (
	"github.com/liu-jiangyuan/go_websocket/conf"
	"github.com/liu-jiangyuan/go_websocket/engine"
	"log"
)

func main()  {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("err:%+v",err)
		}
	}()
	e := engine.InitEngine()
	conf.InitRoute()
	e.SetHost("0.0.0.0")
	e.SetPort("8089")

	//路由设置 可在此处设置，也可以在conf/router.go配置文件设置
	//e.AddRoute("send",controller.Send)
	//e.SetHandle("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	lib.RunServer(w,r)
	//})
	e.Run()
}