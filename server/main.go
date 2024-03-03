package main

import (
	"chatroom/server/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO: enforce limits
var ROOM_CAPACITY = 8
var SERVER_CAPACITY = 32
var MAX_ROOMS = 4

var connectingClients = make(chan *models.Client)     // client type bc displayName is received from connection request
var disconnectingClients = make(chan *websocket.Conn) // conn type bc remote addr is key
var clients = make(map[string]*models.Client)         // {remote addr: client} all clients in server
var rooms = make(map[int]*models.Room)                // {int: room} all rooms in server

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func run() {
	for {
		select {
		case client := <-connectingClients:
			clients[client.Conn.RemoteAddr().String()] = client
			log.Printf("(%s) client connected as %s", client.Conn.RemoteAddr().String(), client.DisplayName)
		case conn := <-disconnectingClients:
			addr := conn.RemoteAddr().String()
			log.Printf("(%s) client disconnected as %s", addr, clients[addr].DisplayName)
			delete(clients, addr)
		}
	}
}

func main() {
	go run()
	http.HandleFunc("/", handleConnection)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
