package main

import (
	"chatroom/server/models"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Rooms struct {
	sync.Mutex
	roomMap map[int]*models.Room
}

var connectingClients = make(chan *models.Client)      // client type bc displayName is received from connection request
var disconnectingConns = make(chan *websocket.Conn)    // conn type bc remote addr is key
var clients = make(map[string]*models.Client)          // {remote addr: client} all clients in server
var rooms = Rooms{roomMap: make(map[int]*models.Room)} // {int: room} all rooms in server

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func run() {
	for {
		select {
		case client := <-connectingClients:
			if len(clients) == SERVER_CAPACITY {
				client.WriteText("Server capacity reached. Please try again later.")
				disconnectingConns <- client.GetConn()
				break
			}
			clients[client.GetConn().RemoteAddr().String()] = client
			log.Printf("(%s) client connected as %s", client.GetConn().RemoteAddr().String(), client.GetDisplayName())
		case conn := <-disconnectingConns:
			addr := conn.RemoteAddr().String()
			log.Printf("(%s) client disconnected as %s", addr, clients[addr].GetDisplayName())
			delete(clients, addr)
		}
	}
}

func main() {
	go run()
	http.HandleFunc("/", handleConnection)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
