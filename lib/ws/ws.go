package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/liu-jiangyuan/go_websocket/lib/gateway"
	"github.com/liu-jiangyuan/go_websocket/lib/msg"
	"log"
	"net/http"
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
	Msg *msg.Msg
}

func (c *Info) Close() {
	c.Conn.Close()
	c.mutex.Lock()
	if !c.isClosed {
		c.isClosed = true
		gateway.Gateway.UnbindClient(c.Conn)
		gateway.Gateway.SendToAll([]byte(fmt.Sprintf("%d is closed",c.Uid)))
	}
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

		//读到的信息全部发送消息中心
		select {
		case c.Msg.In <- message:

		}
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
		initMsg := msg.InitMsg()
		client := &Info{
			Conn: conn,
			Uuid: "",Uid:id,
			Msg:  initMsg,
		}
		gateway.Gateway.UidBindClient(id,conn)
		gateway.Gateway.ClientBindUid(conn,id)

		go client.tickerLoop()
		go client.ReadLoop()
		go initMsg.ParseMsg(conn)
	}

}