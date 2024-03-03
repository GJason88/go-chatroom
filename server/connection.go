package main

import (
	"chatroom/server/models"
	"chatroom/utils"
	"fmt"
	"log"
	"net/http"
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
		disconnectingClients <- conn
	}()
	client := models.CreateClient(r.URL.Query()["displayName"][0], conn)
	connectingClients <- client
	listen(client)
}

func listen(client *models.Client) {
	for {
		_, msgBytes, err := client.Conn.ReadMessage()
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
			listRoomsController(client, rooms)
		case "join":
			if len(args) < 2 {
				client.WriteText("Missing room number.")
				break
			}
			joinRoomController(client, args[1])
		case "create":
			if len(args) < 3 {
				client.WriteText("Missing room name and/or size.")
				break
			}
			createRoomController(client, args[1], args[2])
		case "help":
			client.Help()
		case "quit", "exit":
			return
		default:
			client.WriteText(fmt.Sprintf("Unknown command: %s\nType \"help\" to see all commands.", args[0]))
		}
	}
}
