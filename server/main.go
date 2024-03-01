package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	displayName string
	conn        *websocket.Conn
}

type Room struct {
	connectingUsers    chan *Client
	disconnectingUsers chan *Client
	users              map[string]*Client // addr to client
	messages           chan string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (room Room) connectClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	displayName := r.URL.Query()["displayName"][0]
	log.Printf("(%s) client connected as %s", conn.RemoteAddr().String(), displayName)
	client := Client{
		displayName,
		conn,
	}
	room.connectingUsers <- &client
}

func (room Room) disconnectClient(client *Client) {
	if err := client.conn.Close(); err != nil {
		log.Println("close error:", err)
	}
	room.disconnectingUsers <- client
}

func (room Room) readFromClient(client *Client) {
	defer room.disconnectClient(client)
	for {
		_, msgBytes, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, 1000) {
				log.Printf("(%s) client disconnected as %s", client.conn.RemoteAddr().String(), client.displayName)
			} else {
				log.Println("read error:", err)
			}
			break
		}
		msg := fmt.Sprintf("%s: %s", client.displayName, string(msgBytes))
		log.Printf("(%s) %s", client.conn.RemoteAddr().String(), msg)
		room.messages <- msg
	}
}

func (room Room) writeToClients(msg string) {
	for _, client := range room.users {
		client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func (room Room) run() {
	for {
		select {
		case client := <-room.connectingUsers:
			room.users[client.conn.RemoteAddr().String()] = client
			room.writeToClients(fmt.Sprintf("%s has joined the room", client.displayName))
			go room.readFromClient(client)
		case client := <-room.disconnectingUsers:
			room.writeToClients(fmt.Sprintf("%s has left the room", client.displayName))
			delete(room.users, client.conn.RemoteAddr().String())
		case msg := <-room.messages:
			room.writeToClients(msg)
		}
	}
}

func createRoom() *Room {
	return &Room{
		make(chan *Client),
		make(chan *Client),
		make(map[string]*Client),
		make(chan string),
	}
}

func main() {
	room := createRoom()
	go room.run()
	http.HandleFunc("/join", room.connectClient)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
