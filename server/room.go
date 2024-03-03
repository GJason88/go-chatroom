package main

type Room struct {
	number               int
	name                 string
	connectingClients    chan *Client
	disconnectingClients chan *Client
	clients              map[string]*Client
	messages             chan string
}

// func (room Room) readFromClient(client *Client) {
// 	defer disconnectClient(client)
// 	for {
// 		_, msgBytes, err := client.conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsCloseError(err, 1000) {
// 				log.Printf("(%s) client disconnected as %s", client.conn.RemoteAddr().String(), client.displayName)
// 			} else {
// 				log.Println("read error:", err)
// 			}
// 			break
// 		}
// 		msg := fmt.Sprintf("%s: %s", client.displayName, string(msgBytes))
// 		log.Printf("(%s) %s", client.conn.RemoteAddr().String(), msg)
// 		room.messages <- msg
// 	}
// }

// func (room Room) writeToClients(msg string) {
// 	for _, client := range room.users {
// 		client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
// 	}
// }

func createRoom() *Room {
	autoIncRoomId.Lock()
	defer autoIncRoomId.Unlock()
	autoIncRoomId.id++
	return &Room{
		autoIncRoomId.id,
		"room",
		make(chan *Client),
		make(chan *Client),
		make(map[string]*Client),
		make(chan string),
	}
}

// room.writeToClients(fmt.Sprintf("%s has joined the room", client.displayName))
// go room.readFromClient(client)
// room.writeToClients(fmt.Sprintf("%s has left the room", client.displayName))
