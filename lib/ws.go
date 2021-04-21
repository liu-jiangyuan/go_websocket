package lib

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/liu-jiangyuan/go_websocket/conf"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Info struct {
	Conn *websocket.Conn
	Uuid string
	Uid int64
	mutex sync.Mutex  // 对closeChan关闭上锁
	isClosed bool  // 防止closeChan被关闭多次
}
type msg struct {
	Method    string `json:"method"`
	Data      map[string]interface{} `json:"data"`
}


func (c *Info) Close() {
	c.Conn.Close()
	c.mutex.Lock()
	if !c.isClosed {
		c.isClosed = true
	}
	Gateway.UnbindClient(c.Conn)
	Gateway.SendToAll([]byte(fmt.Sprintf("%d is closed",c.Uid)))
	c.mutex.Unlock()
}

func (c *Info) ReadLoop () {
	defer func() {
		if err := recover(); err != nil {
			c.Close()
			log.Printf("fatal error:%+v",err)
		}
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	},)

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Close()
				log.Printf("err:%+v",err)
			}
			break
		}

		var parse msg
		if string(message) == "PING" {
			parse = msg{
				Method: "PING",
				Data:   map[string]interface{}{"client":c},
			}
		} else {
			err = json.Unmarshal(message,&parse)
			if err != nil {
				panic(err)
			}
		}

		if _, ok := conf.RouteMap[parse.Method]; !ok {
			c.Conn.WriteMessage(websocket.TextMessage,[]byte("method not Exits"))
			break
		}
		fv := reflect.ValueOf(conf.RouteMap[parse.Method])
		params := make([]reflect.Value,1)  //参数
		params[0] = reflect.ValueOf(parse.Data)
		res := fv.Call(params)
		a := res[0].Interface().(map[string]interface{})

		r , _ := json.Marshal(a)
		c.Conn.WriteMessage(websocket.TextMessage,r)
		//Gateway.SendToAll(message)
		//log.Printf("ReadLoop message:%+v",string(res[0].Interface().([]byte)))
	}
}
func (c *Info) tickerLoop ()  {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil { //心跳检测失败，默认离线
				c.Close()
				return
			}
		}
	}
}

func RunServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	param := r.URL.Query()
	if id , err := strconv.ParseInt(param["id"][0],10,64); err == nil {
		client := &Info{
			Conn: conn,
			Uuid:"",Uid:id,
		}
		Gateway.UidBindClient(id,client)
		Gateway.ClientBindUid(conn,client)

		go client.tickerLoop()
		go client.ReadLoop()
	}

}