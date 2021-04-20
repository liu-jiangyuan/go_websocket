package lib

import (
	"github.com/gorilla/websocket"
)

type gateWay struct {
	UidBindClientMap map[int64]*Info
	ClientBindUidMap map[*websocket.Conn]*Info
	GroupMap map[string][]int64
}

var Gateway *gateWay

func init() {
	Gateway = &gateWay{
		UidBindClientMap: make(map[int64]*Info),
		ClientBindUidMap: make(map[*websocket.Conn]*Info),
		GroupMap:make(map[string][]int64),
	}
}

func Asdf(msg []byte) {
	Gateway.SendToAll(msg)
}

func (g *gateWay) UidBindClient (uid int64,client *Info) {
	g.UidBindClientMap[uid] = client
}

func (g *gateWay) ClientBindUid (conn *websocket.Conn,client *Info) {
	g.ClientBindUidMap[conn] = client
}

func (g *gateWay) UnbindUid (uid int64) {
	delete(g.UidBindClientMap,uid)
}

func (g *gateWay) UnbindClient (conn *websocket.Conn) {
	delete(g.ClientBindUidMap,conn)
}

func (g *gateWay) JoinGroup (uid int64,groupName string) {
	g.GroupMap[groupName] = append(g.GroupMap[groupName],uid)
}

func (g *gateWay) LeaveGroup (uid int64,groupName string) {

}

func (g *gateWay) SendToUid (uid int64,msg []byte) error {
	return g.UidBindClientMap[uid].Conn.WriteMessage(websocket.TextMessage,msg)
}

func (g *gateWay) SendToGroup (groupName string,msg []byte) {
	for _, uid := range g.GroupMap[groupName] {
		g.UidBindClientMap[uid].Conn.WriteMessage(websocket.TextMessage,msg)
	}
}

func (g *gateWay) SendToAll (msg []byte) {
	for conn, _ := range g.ClientBindUidMap {
		conn.WriteMessage(websocket.TextMessage,msg)
	}
}