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
	roomSize, err := strconv.Atoi(roomSizeStr)
	if err != nil {
		client.WriteText("Please enter a valid number for room size.")
		return
	}
	if MIN_ROOM_SIZE > roomSize || roomSize > MAX_ROOM_SIZE {
		client.WriteText(fmt.Sprintf("Please enter a room size between %d and %d.", MIN_ROOM_SIZE, MAX_ROOM_SIZE))
		return
	}
	if len(rooms) == MAX_ROOMS {
		client.WriteText("Max number of rooms reached. Please join an existing room or try again later.")
	}
	// room := models.CreateRoom(roomName, roomSize)
}

func listRoomsController(client *models.Client, rooms map[int]*models.Room) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Room Number\tRoom Name\tUsers")
	for _, room := range rooms {
		fmt.Fprintf(w, "%d\t%s\t%d/%d\n", room.Number, room.Name, len(room.Clients), room.Size)
	}
	client.WriteText(buf.String())
}
