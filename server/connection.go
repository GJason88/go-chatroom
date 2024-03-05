package main

import (
	"chatroom/server/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connectClient(w http.ResponseWriter, r *http.Request) {
	clients.Lock()
	defer clients.Unlock()
	if len(clients.clientMap) == SERVER_CAPACITY {
		w.Write([]byte("Server capacity reached. Please try again later."))
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	addr := conn.RemoteAddr().String()
	displayName := r.URL.Query()["displayName"][0]

	client := models.CreateClient(displayName, conn)
	clients.clientMap[addr] = client

	log.Printf("(%s) client connected as %s", addr, displayName)
	go listen(client)
}

func disconnectClient(client *models.Client) {
	addr, displayName := client.Disconnect()
	log.Printf("(%s) client disconnected as %s", addr, displayName)
	delete(clients.clientMap, addr)
}

// TODO: figure out how to stop listening to client after they join a room
// TODO: delete rooms when no more clients in them
func listen(client *models.Client) {
	defer disconnectClient(client)
	for {
		msg, err := client.ReadText()
		if err != nil {
			break
		}
		args := strings.Fields(msg)
		if len(args) == 0 {
			continue
		}
		switch args[0] {
		case "rooms":
			listRooms(client, rooms.roomMap)
		case "join":
			if len(args) < 2 {
				client.WriteText("Missing room number.")
				break
			}
			roomNumber, err := strconv.Atoi(args[1])
			if err != nil {
				client.WriteText("Please enter a valid number.")
				return
			}
			addClientToRoom(client, roomNumber)
		case "create":
			if len(args) < 3 {
				client.WriteText("Missing room name and/or size.")
				break
			}
			if room := createRoom(client, args[1], args[2]); room != nil {
				go runRoom(room)
				addClientToRoom(client, room.GetNumber())
			}
		case "help":
			client.Help()
		case "quit", "exit":
			return
		default:
			client.WriteText(fmt.Sprintf("Unknown command: %s\nType \"help\" to see all commands.", args[0]))
		}
	}
}
