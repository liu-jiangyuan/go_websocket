package lib

import (
	"github.com/gorilla/websocket"
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

var upgrade = websocket.Upgrader{
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
	Msg *Msg
}

func (c *Info) close() {
	c.Conn.Close()
	c.mutex.Lock()
	if !c.isClosed {
		c.isClosed = true
		c.Msg.UnBind <- [2]interface{}{c.Uid,c.Conn}
	}
	c.mutex.Unlock()
}

func (c *Info) readLoop () {
	defer func() {
		if err := recover(); err != nil {
			c.close()
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
				c.close()
				break
			}
			log.Printf("err:%+v",err)
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
				c.close()
				return
			}
		}
	}
}

func RunServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	param := r.URL.Query()
	if id , err := strconv.ParseInt(param["id"][0],10,64); err == nil {
		initMsg := InitMsg()
		client := &Info{
			Conn: conn,
			Uuid: "",Uid:id,
			Msg:  initMsg,
		}

		go client.tickerLoop()
		go client.readLoop()
		go initMsg.Parse(conn)

		//登录绑定
		initMsg.Bind <- [2]interface{}{id,conn}
	}
}