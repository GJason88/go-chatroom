package models

import (
	"fmt"
	"sync"
)

type Clients struct {
	sync.Mutex
	clientMap map[string]*Client
}

type Room struct {
	number   int
	name     string
	capacity int
	clients  map[string]*Client
	messages chan string
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

func (r *Room) writeToClients(msg string) {
	for _, c := range r.clients {
		c.WriteText(msg)
	}
}

func (r *Room) AddAndListen(client *Client) {
	for {
		if len(r.clients) == 0 {
			break
		}
		msg, err := client.ReadText()
		if err != nil || msg == "/leave" {
			break
		}
		r.messages <- fmt.Sprintf("%s: %s", client.GetDisplayName(), msg)
	}
}

func (r *Room) removeClient(client *Client) {

}

func (r *Room) GetCapacity() int {
	return r.capacity
}

func (r *Room) GetNumber() int {
	return r.number
}

func (r *Room) GetName() string {
	return r.name
}

func (r *Room) GetHeadCount() int {
	return len(r.clients)
}

func (r *Room) Run() {
	// for {
	// 	if len(r.clients) == 0 {
	// 		break
	// 	}
	// 	select {
	// 	case client := <-r.joiningClients:
	// 		if len(r.clients) == r.capacity {
	// 			client.WriteText("Room is full.")
	// 			r.leavingClients <- client
	// 			break
	// 		}
	// 		r.clients[client.GetConn().RemoteAddr().String()] = client
	// 		client.WriteText("You have connected to the room. Type \"/leave\" to leave.")
	// 		r.writeToClients(fmt.Sprintf("%s has joined the room", client.displayName))
	// 	case client := <-r.leavingClients:
	// 		// TODO:
	// 		fmt.Println(client.GetDisplayName())
	// 	case msg := <-r.messages:
	// 		// TODO:
	// 		fmt.Println(msg)
	// 	}
	// }
}

func CreateRoom(roomName string, size int) *Room {
	autoIncRoomId.Lock()
	defer autoIncRoomId.Unlock()
	autoIncRoomId.id++
	return &Room{
		autoIncRoomId.id,
		roomName,
		size,
		make(chan *Client),
		make(chan *Client),
		make(map[string]*Client),
		make(chan string),
	}
}

// go room.readFromClient(client)
// room.writeToClients(fmt.Sprintf("%s has left the room", client.displayName))
