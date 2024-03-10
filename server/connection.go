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
	if len(clients.clientMap) == SERVER_CAPACITY {
		clients.Unlock()
		w.Write([]byte("Server capacity reached. Please try again later."))
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		clients.Unlock()
		log.Println("upgrade error:", err)
		return
	}

	addr := conn.RemoteAddr().String()
	displayName := r.URL.Query()["displayName"][0]
	log.Printf("(%s) client connected as %s", addr, displayName)

	client := models.CreateClient(displayName, conn)
	defer disconnectClient(client)
	clients.clientMap[addr] = client
	clients.Unlock()
	listen(client)
}

func listen(client *models.Client) {
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
				break
			}
			room, ok := rooms.roomMap[roomNumber]
			if !ok {
				client.WriteText("Room does not exist.")
				break
			}
			room.AddClient(client, false) // blocking
		case "create":
			if len(args) < 3 {
				client.WriteText("Missing room name and/or size.")
				break
			}
			if room := createRoom(client, args[1], args[2]); room != nil {
				go runRoom(room)
				room.AddClient(client, true) // blocking
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

func disconnectClient(client *models.Client) {
	addr, displayName := client.Disconnect()
	log.Printf("(%s) client disconnected as %s", addr, displayName)
	clients.Lock()
	defer clients.Unlock()
	delete(clients.clientMap, addr)
}
