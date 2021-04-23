package test

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func server() {
	log.SetFlags(log.LstdFlags)
	log.Println(time.Now().UnixNano())
	var wg sync.WaitGroup
	userMap := 0
	reviceMap := 0

	origin := "http://localhost/"
	var rev = make(chan []byte,2048)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				//log.Printf("fatal error:%+v",err)
			}
		}()
		for {
			select {
			case <- rev:
				//fmt.Printf("go func: %s.\n", data)
			}
		}
	}()
	for id := 1;id <= 3000;id ++  {
		wg.Add(1)
		userMap ++
		go func(id int) {
			defer func() {
				if err := recover(); err != nil {
					//log.Printf("fatal error:%+v",err)
				}
			}()
			url := fmt.Sprintf("ws://localhost:8089/ws?id=%d",id)
			ws, err := websocket.Dial(url, "", origin)
			if err != nil {
				panic(err)
			}

			if _, err := ws.Write([]byte("PING")); err != nil {
				panic(err)
			}
			go func(id int) {
				wg.Add(1)
				var msg = make([]byte, 1024)
				var n int
				if n, err = ws.Read(msg); err != nil {
					panic(err)
				} else {
					reviceMap ++
					rev <- msg[:n]
				}
				wg.Done()
			}(id)
			wg.Done()
		}(id)
	}

	//time.Sleep(time.Second * 1)
	wg.Wait()
	log.Printf("send:%d;revice:%+d",userMap,reviceMap)
	log.Println(time.Now().UnixNano())
}
