package main

import (
	"fmt"
	"strconv"
	"strings"
)

//check functions mei room reference se pass nahi kiya

var room_list []room

func init() {
	room_list = []room{}
}

type client struct {
	name string
	ws   websocket
}
type room struct {
	clients      []client
	broadcast    chan string
	sender       int
	current_word string
	token        string
	rounds_left  int
}

func newroom(ws websocket, user, tk string) {
	temp := make([]byte, 100)
	ws.read(temp)
	rounds_left := strconv.Atoi(string(temp))
	r := room{
		clients:      []client{client{name: user, ws: ws}},
		broadcast:    make(chan string),
		sender:       0,
		current_word: "",
		token:        tk,
		rounds_left:  rounds_left,
	}
	go game_loop(&r)
	return r
}

func selectsyncdb(ws websocket, user string, token string) {
	for i, r := range room_list {
		if r.token == token {
			r.clients = append(r.clients, client{name: user, ws: ws})
			return &r, i
		}
	}
	return newroom(ws, user, token), 0
}

func left(i int, r room) {
	room.clients = append(room.clients[:i], room.clients[i+1:]...)
}

func wshandl(ws websocket) {
	res1 := strings.Split(ws.LocalAddr().String(), "/")
	if len(res1) < 5 {
		ws.Close()
	}

	token, err := strconv.Atoi(res1[5])
	if err != nil {
		fmt.Println("closed because of ", err)
		ws.Close()
		return
	}

	temp := make([]byte, 100)
	_, err = ws.Read(temp)
	if err != nil {
		fmt.Println("closed because of ", err)
		ws.Close()
		return
	}

	user := string(temp)
	fmt.Println(token, user)

	r, i := selectsyncdb(ws, user, token)

	for {
		if i != r.sender {
			temp := make([]byte, 100)
			ws.Read(temp)
			op := string(temp)
			if op == "p_1" { //point plus 1
				broadcast("p_1_"+user, r)
			}
			if op == "close" { //point plus 1
				left(i, r)
				break
			}
		}
	}
	ws.close()

	//sync waitgroup

}
