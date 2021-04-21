package msg

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/liu-jiangyuan/go_websocket/engine"
	"log"
	"reflect"
)

type Msg struct {
	In chan []byte
	Close chan []byte
}

type parse struct {
	Method    string `json:"method"`
	Data      map[string]interface{} `json:"data"`
}

func InitMsg() *Msg {
	return &Msg{
		In:make(chan []byte),
		Close:make(chan []byte),
	}
}

func (m *Msg) ParseMsg(conn *websocket.Conn) {
	var(
		data []byte
		err error
	)
	for {
		select{
		case data = <- m.In:
		}
		var p parse
		if string(data) == "PING" {
			p = parse{
				Method: "PING",
				Data:   map[string]interface{}{"PONG":"PONG"},
			}
		} else {
			err = json.Unmarshal(data,&p)
			if err != nil {
				panic(err)
			}
		}

		route := engine.GetRoute()
		if _, ok := route[p.Method]; !ok {
			log.Println("method not Exits")
			//c.Conn.WriteMessage(websocket.TextMessage,[]byte("method not Exits"))
			break
		}
		fv := reflect.ValueOf(route[p.Method])
		params := make([]reflect.Value,1)  //参数
		params[0] = reflect.ValueOf(p.Data)
		//fv.Call(params)
		res := fv.Call(params)
		a := res[0].Interface().(map[string]interface{})

		r , _ := json.Marshal(a)
		conn.WriteMessage(websocket.TextMessage,r)
		//Gateway.SendToAll(message)
		//log.Printf("ReadLoop message:%+v",string(res[0].Interface().([]byte)))
	}
}