package models

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Room struct {
	number         int
	name           string
	capacity       int
	joiningClients chan *Client
	leavingClients chan *Client
	clients        map[string]*Client
	messages       chan string
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
		c.GetConn().WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func (r *Room) AddClient(client *Client) {
	r.joiningClients <- client
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
	for {
		if len(r.clients) == 0 {
			break
		}
		select {
		case client := <-r.joiningClients:
			if len(r.clients) == r.capacity {
				client.WriteText("Room is full.")
				r.leavingClients <- client
				break
			}
			r.writeToClients(fmt.Sprintf("%s has joined the room", client.displayName))
			r.clients[client.GetConn().RemoteAddr().String()] = client
			// TODO: goroutine to read from client? Need to find solution to stop reading from main connection thread.
		case client := <-r.leavingClients:
			// TODO:
			fmt.Println(client.GetDisplayName())
		case msg := <-r.messages:
			// TODO:
			fmt.Println(msg)
		}
	}
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
