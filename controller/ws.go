package controller

import (
	"encoding/json"
	"github.com/liu-jiangyuan/go_websocket/gateway"
)

func Index (param map[string]interface{}) map[string]interface{} {
	b , _:= json.Marshal(param)
	gateway.Gateway.SendToAll(b)
	return param
}


func Send (param map[string]interface{}) map[string]interface{} {
	b , _:= json.Marshal(param)
	gateway.Gateway.SendToAll(b)
	return param
}