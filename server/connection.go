package main

import (
	"chatroom/server/models"
	"chatroom/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer func() {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
		disconnectingConns <- conn
	}()
	client := models.CreateClient(r.URL.Query()["displayName"][0], conn)
	connectingClients <- client
	listen(client)
}

// TODO: figure out how to stop listening to client after they join a room
// TODO: delete rooms when no more clients in them
func listen(client *models.Client) {
	for {
		_, msgBytes, err := client.GetConn().ReadMessage()
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
