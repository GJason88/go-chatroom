package main

import (
	"bytes"
	"chatroom/server/models"
	"fmt"
	"strconv"
	"text/tabwriter"
)

func joinRoomController(client *models.Client, roomNumberStr string) {
	roomNumber, err := strconv.Atoi(roomNumberStr)
	if err != nil {
		client.WriteText("Please enter a valid number.")
		return
	}
	room, ok := rooms[roomNumber]
	if !ok {
		client.WriteText("Room does not exist.")
		return
	}
	room.AddClient(client)
}

func createRoomController(client *models.Client, roomName string, roomSizeStr string) {
	client.WriteText("TODO: create room")
	// roomSize, err := strconv.Atoi(roomSizeStr)
	// if err != nil {
	// 	client.WriteText("Please enter a valid number for room size.")
	// 	return
	// }
	// room := models.CreateRoom(roomName, roomSize)
}

func listRoomsController(client *models.Client, rooms map[int]*models.Room) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Room Number\tRoom Name\tUsers")
	for _, room := range rooms {
		fmt.Fprintf(w, "%v\t%s\t%v/%v\n", room.Number, room.Name, len(room.Clients), room.Size)
	}
	client.WriteText(buf.String())
}
