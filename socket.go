package main

import (
	"github.com/liu-jiangyuan/go_websocket/conf"
	"github.com/liu-jiangyuan/go_websocket/controller"
	"github.com/liu-jiangyuan/go_websocket/engine"
	"github.com/liu-jiangyuan/go_websocket/lib"
	"log"
	"time"
)

func main()  {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	e := engine.InitEngine()
	conf.InitRoute()
	e.SetHost("0.0.0.0")
	e.SetPort("8089")
	id := lib.Timer.After(lib.TimerStruct{
		D:      7 * time.Second,
		Args: map[string]interface{}{"method":"AfterWithFuncArgsNoClear"},
		Action: controller.Index,
	})
	log.Printf("lib.Timer After id:%+v",id)
	id2 := lib.Timer.After(lib.TimerStruct{
		D:      9 * time.Second,
		Args: map[string]interface{}{"method":"AfterWithFuncArgs"},
		Action: controller.Index,
	})
	log.Printf("lib.Timer Tick id:%+v",id2)

	id3 := lib.Timer.Loop (lib.TimerStruct{
		D:      1 * time.Second,
		Args: map[string]interface{}{"method":"LoopWithFuncArgs"},
		Action: controller.Index,
	},1,5)
	log.Printf("lib.Timer Tick id:%+v",id3)

	lib.Timer.MustOnceAfter(lib.TimerStruct{
		D:      15 * time.Second,
		Args: map[string]interface{}{"method":"MustOnceAfterWithFunArgs"},
		Action: controller.Index,
	})

	time.AfterFunc(3 * time.Second, func() {
		lib.Timer.Clear(id3)
	})
	lib.Timer.Clear(id3)

	//路由设置 可在此处设置，也可以在conf/router.go配置文件设置
	//e.AddRoute("send",controller.Send)
	//e.SetHandle("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	lib.RunServer(w,r)
	//})
	e.Run()
}