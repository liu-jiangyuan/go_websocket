package test

import (
	"github.com/liu-jiangyuan/go_websocket/controller"
	"github.com/liu-jiangyuan/go_websocket/lib"
	"log"
	"time"
)

func Call (args ...map[string]interface{}) map[string]interface{} {
	//log.Printf("-------------- Call args:%+v --------------------",args)
	r := map[string]interface{}{"A":"B"}
	if len(args) > 0 {
		r = args[0]
	}
	return r
}

func timer() {
	id := lib.Timer.After(time.Second * 7 , Call)
	log.Printf("lib.Timer After id:%+v",id)

	id2 := lib.Timer.After(time.Second * 3 , Call)
	log.Printf("lib.Timer Tick id:%+v",id2)

	id3 := lib.Timer.Loop (time.Second * 1 , func(args ...map[string]interface{}) map[string]interface{} {
		return Call(map[string]interface{}{"method":"loop-call","test1":1})
	},1,5)
	log.Printf("lib.Timer Tick id:%+v",id3)

	lib.Timer.MustOnceAfter(time.Second * 5 , func(args ...map[string]interface{}) map[string]interface{} {
		return Call(controller.Index(map[string]interface{}{"method":"MustOnceAfter"}))
	})

	time.AfterFunc(3 * time.Second, func() {
		lib.Timer.Clear(id3)
	})
	lib.Timer.Clear(id)
}