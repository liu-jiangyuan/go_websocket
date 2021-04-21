package controller

import (
	"encoding/json"
	"github.com/liu-jiangyuan/go_websocket/lib/gateway"
	"log"
	"net/http"
)

func Index (param map[string]interface{}) map[string]interface{} {
	b , _:= json.Marshal(param)
	gateway.Gateway.SendToAll(b)
	return param
}

func Test(w http.ResponseWriter, r *http.Request)  {
	_, _ = w.Write([]byte("there is Test"))
}

func Aa (param map[string]interface{}) map[string]interface{} {
	log.Printf("%+v",param)
	return param
}