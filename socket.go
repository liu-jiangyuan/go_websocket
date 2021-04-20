package main

import (
	"github.com/liu-jiangyuan/go_websocket/engine"
	"log"
)

func main()  {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	e := engine.InitEngine()
	e.Host = "0.0.0.0"
	e.Port = "8089"
	e.Run()
}
