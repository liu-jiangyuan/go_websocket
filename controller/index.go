package controller

import (
	"log"
	"net/http"
)

func Index (param map[string]interface{}) map[string]interface{} {
	return param
}

func Test(w http.ResponseWriter, r *http.Request)  {
	_, _ = w.Write([]byte("there is Test"))
}

func Aa (param map[string]interface{}) map[string]interface{} {
	log.Printf("%+v",param)
	return param
}