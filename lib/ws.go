package lib

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/liu-jiangyuan/go_websocket/conf"
	"log"
	"net/http"
	"reflect"
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
var ClientPool []*Client

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


type Client struct {
	Conn *websocket.Conn
	Uuid string
	Uid int64
}
type msg struct {
	Method    string `json:"method"`
	Data      map[string]interface{} `json:"data"`
}
func (c *Client) ReadLoop () {
	defer func() {
		log.Println("fatal error")
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
				log.Printf("err:%+v",err)
			}
			break
		}

		var parse msg
		err = json.Unmarshal(message,&parse)
		if err != nil {
			log.Fatal(err)
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
		//log.Printf("ReadLoop message:%+v",string(res[0].Interface().([]byte)))
	}
}

func (s *Client) tickerLoop ()  {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case <-ticker.C:
			if err := s.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {//心跳检测失败，默认离线
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

	client := &Client{Conn: conn,Uuid:"",Uid:0}
	ClientPool = append(ClientPool,client)

	go client.tickerLoop()
	go client.ReadLoop()
}