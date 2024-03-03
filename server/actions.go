package main

import (
	"strconv"
)

func roomsAction(client *Client) {
	client.writeText("TODO: list rooms action")
}

func joinAction(msg string, client *Client) {
	roomNumber, err := strconv.Atoi(msg)
	if err != nil {
		client.writeText("Please enter a valid number.")
		return
	}
	room, ok := rooms[roomNumber]
	if !ok {
		client.writeText("Room does not exist.")
		return
	}
	room.connectingClients <- client
	// room.listen(client) // block and listen in room until client disconnects
}

func createAction(roomName string, client *Client) {
	client.writeText("TODO: create action")
}

func helpAction(client *Client) {
	client.writeText("TODO: Help action")
}

func quitAction(client *Client) {
	client.writeText(("TODO: Quit action"))
}
