package lib

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/liu-jiangyuan/go_websocket/gateway"
	"log"
	"net/http"
	"reflect"
)

//route相关，websocket动态调用,所以配置在这里
var Route route
type route struct {
	Ws   map[string]func(map[string]interface{}) map[string]interface{}
	Http map[string]func(w http.ResponseWriter, r *http.Request)
}

func init() {
	Route = route{
		Ws: make(map[string]func(map[string]interface{}) map[string]interface{}),
		Http: make(map[string]func(w http.ResponseWriter, r *http.Request)),
	}
}


//消息channel
type Msg struct {
	In      chan []byte
	Close   chan []byte
	Bind    chan [2]interface{}
	UnBind  chan [2]interface{}
}

func InitMsg() *Msg {
	return &Msg{
		In:     make(chan []byte),
		Close:  make(chan []byte),
		Bind:   make(chan [2]interface{}),
		UnBind: make(chan [2]interface{}),
	}
}

func (m *Msg) Parse(conn *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			//log.Printf("fatal error:%+v",err)
		}
	}()
	var (
		in []byte
		bind [2]interface{}
		unbind [2]interface{}
	)
	for {
		select {
		case in = <- m.In:
			m.controller(conn,in)
			break
		case bind = <- m.Bind:
			m.bind(bind)
			break
		case unbind = <- m.UnBind:
			m.unBind(unbind)
		}
	}
}

//websocket动态调用方法
type parse struct {
	Method    string `json:"method"`
	Data      map[string]interface{} `json:"data"`
}
func (m *Msg) controller (conn *websocket.Conn,data []byte) {
	var p parse
	if string(data) == "PING" {
		p = parse{
			Method: "PING",
			Data:   map[string]interface{}{"PONG":"PONG"},
		}
	} else {
		err := json.Unmarshal(data,&p)
		if err != nil {
			panic(err)
		}
	}

	route := Route.Ws
	if _, ok := route[p.Method]; !ok {
		log.Println("method not Exits")
		return
	}
	fv := reflect.ValueOf(route[p.Method])
	params := make([]reflect.Value,1)  //参数
	params[0] = reflect.ValueOf(p.Data)
	//fv.Call(params)
	res := fv.Call(params)
	a := res[0].Interface().(map[string]interface{})

	r , _ := json.Marshal(a)
	conn.WriteMessage(websocket.TextMessage,r)
	return
}

//实例话struct方式调用方法，未完成
//type Handler struct {
//	Func  reflect.Value
//	In   reflect.Type
//	NumIn int
//	Out  reflect.Type
//	NumOut int
//}
//func InitRouter() {
//handlers := make(map[string]*lib.Handler)
//v := reflect.ValueOf(&controller.Index{})
//t := reflect.TypeOf(&controller.Index{})
//for i := 0; i < v.NumMethod(); i++ {
//	name := t.Method(i).Name
//	// 可以根据 i 来获取实例的方法，也可以用 v.MethodByName(name) 获取
//	m := v.Method(i)
//	// 这个例子我们只获取第一个输入参数和第一个返回参数
//	in := m.Type().In(0)
//	out := m.Type().Out(0)
//	handlers[name] = &lib.Handler{
//		Func:  m,
//		In:   in,
//		NumIn: m.Type().NumIn(),
//		Out:  out,
//		NumOut: m.Type().NumOut(),
//	}
//}
//return handlers
//inVal := reflect.New(r[c].In).Elem()
//rtn := r[c].Func.Call([]reflect.Value{inVal})[0]
//}


func (m *Msg) bind (data [2]interface{}) {
	uid := data[0].(int64)
	conn := data[1].(*websocket.Conn)
	gateway.Gateway.UidBindClient(uid,conn)
	gateway.Gateway.ClientBindUid(conn,uid)
	return
}

func (m *Msg) unBind (data [2]interface{}) {
	uid := data[0].(int64)
	conn := data[1].(*websocket.Conn)
	gateway.Gateway.UnbindUid(uid)
	gateway.Gateway.UnbindClient(conn)
	gateway.Gateway.SendToAll([]byte(fmt.Sprintf("%+v is Close",uid)))
	return
}