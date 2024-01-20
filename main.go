package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {

	mex := http.NewServeMux()
	mex.Handle("/game_com", websocket.Handlera(wshandl))
	http.ListenAndServe(":8000", mex)
}
