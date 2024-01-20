package main

func broadcast(msg string, r room) {
	for _, c := range r.clients {
		c.ws.Write([]byte(msg))
	}
}

// func guessing(room){
// 	for i:=0;;i=(i+1)%len(room.clients) {
// 		temp := make([]byte, 100)
// 		select{
// 			case room.clients[i].ws.Read(temp) {

// 			}
// 		}
// 	}

// }

func game_loop(r room) {
	for {
		if room.rounds_left <= 0 {
			broadcast("round_end", r)
		}
		sender := room.clients[room.sender]
		op := ""
		temp := make([]byte, 100)

		sender.ws.Write([]byte("u_sender"))

		for op != "start" {
			sender.ws.Read(temp)
			op = string(temp)
		}

		broadcast("gs_"+sender.name, r) //game_start
		sender.ws.Write([]byte(words))

		sender.ws.Read(temp)
		op = string(temp)
		broadcast("gw_"+op, r) //game_word

		//go guessing() word right check on client jsut send got right and wrong

		for op != "draw_end" {
			sender.ws.Read(temp)
			op = string(temp)
			broadcast(op, r)
			if op == "close" {
				// if sender closes
				left(room.sender, r)
				break
			}
		}

		room.sender = (room.sender + 1) % len(room.clients)
		room.rounds_left -= 1
	}

}
