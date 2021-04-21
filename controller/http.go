package controller

import "net/http"

func Test(w http.ResponseWriter, r *http.Request)  {
	_, _ = w.Write([]byte("there is Test"))
}
