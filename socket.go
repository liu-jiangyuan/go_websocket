package main

import (
	"github.com/liu-jiangyuan/go_websocket/lib"
	"log"
)

func main()  {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	e := lib.InitEngine()
	//fv := reflect.ValueOf(controller.Aa)
	//params := make([]reflect.Value,1)  //参数
	//params[0] = reflect.ValueOf(map[string]interface{}{"a":"b"})
	//r := fv.Call(params)
	//log.Printf("%+v",r[0].Interface().(map[string]interface{}))
	e.Host = "0.0.0.0"
	e.Port = "8089"
	e.Run()
}
