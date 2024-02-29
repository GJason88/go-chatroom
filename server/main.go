package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	displayName string
	conn        *websocket.Conn
}

type Room struct {
	connectingUsers    chan Client
	disconnectingUsers chan Client
	users              map[string]Client // addr to client
	messages           chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (room Room) handleClientConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	log.Println("client connected:", conn.RemoteAddr().String())
	client := Client{
		r.URL.Query()["displayName"][0],
		conn,
	}
	room.connectingUsers <- client
}

func (room Room) handleClientDisconnect(client Client) {
	log.Println("disconnecting client:", client.conn.RemoteAddr().String())
	if err := client.conn.Close(); err != nil {
		log.Println("close:", err)
	}
	room.disconnectingUsers <- client
}

func (room Room) handleClientReads(client Client) {
	defer room.handleClientDisconnect(client)
	for {
		_, msgBytes, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Println("message:", string(msgBytes))
		room.messages <- msgBytes
	}
}

func (room Room) run() {
	for {
		select {
		case client := <-room.connectingUsers:
			room.users[client.conn.RemoteAddr().String()] = client
			log.Println("client", client.conn.RemoteAddr().String(), "joined room")
			go room.handleClientReads(client)
		case client := <-room.disconnectingUsers:
			delete(room.users, client.conn.RemoteAddr().String())
			log.Println("client", client.conn.RemoteAddr().String(), "left room")
		case msg := <-room.messages:
			for _, client := range room.users {
				client.conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}

func createRoom() *Room {
	return &Room{
		make(chan Client),
		make(chan Client),
		make(map[string]Client),
		make(chan []byte),
	}
}

func main() {
	room := createRoom()
	go room.run()
	http.HandleFunc("/join", room.handleClientConnect)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
