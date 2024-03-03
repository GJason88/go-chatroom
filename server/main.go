package main

import (
	"chatroom/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// TODO: enforce limits
var ROOM_CAPACITY = 8
var SERVER_CAPACITY = 32
var MAX_ROOMS = 4

var connectingClients = make(chan *Client)            // client type bc displayName is received from connection request
var disconnectingClients = make(chan *websocket.Conn) // conn type bc remote addr is key
var clients = make(map[string]*Client)                // {remote addr: client} all clients in server
var rooms = make(map[int]*Room)                       // {int: room} all rooms in server

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func listen(client *Client) {
	for {
		_, msgBytes, err := client.conn.ReadMessage()
		if err != nil {
			utils.LogReadErrors(err)
			break
		}
		args := strings.Fields(string(msgBytes))
		if len(args) == 0 {
			continue
		}
		switch args[0] {
		case "rooms":
			roomsAction(client)
		case "join":
			if len(args) < 2 {
				client.writeText("Missing room number.")
				break
			}
			joinAction(args[1], client)
		case "create":
			if len(args) < 2 {
				client.writeText("Missing room name.")
				break
			}
			createAction(args[1], client)
		case "help":
			helpAction(client)
		case "quit", "exit":
			quitAction(client)
		default:
			client.writeText(fmt.Sprintf("Unknown command: %s\nType \"help\" to see all commands.", args[0]))
		}
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer func() {
		conn.Close()
		disconnectingClients <- conn
	}()
	client := createClient(r.URL.Query()["displayName"][0], conn)
	connectingClients <- client
	listen(client)
}

func run() {
	for {
		select {
		case client := <-connectingClients:
			clients[client.conn.RemoteAddr().String()] = client
			log.Printf("(%s) client connected as %s", client.conn.RemoteAddr().String(), client.displayName)
		case conn := <-disconnectingClients:
			addr := conn.RemoteAddr().String()
			log.Printf("(%s) client disconnected as %s", addr, clients[addr].displayName)
			delete(clients, addr)
		}
	}
}

func main() {
	go run()
	http.HandleFunc("/", handleConnection)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
