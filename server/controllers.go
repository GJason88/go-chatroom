package main

import (
	"bytes"
	"chatroom/server/models"
	"fmt"
	"log"
	"strconv"
	"text/tabwriter"
)

// Adds client to a room
func addClientToRoom(client *models.Client, roomNumber int) {
	room, ok := rooms.roomMap[roomNumber]
	if !ok {
		client.WriteText("Room does not exist.")
		return
	}
	room.AddClient(client)
}

// Creates a room and adds it to the server, returns the room
func createRoom(client *models.Client, roomName string, roomSizeStr string) *models.Room {
	roomSize, err := strconv.Atoi(roomSizeStr)
	if err != nil {
		client.WriteText("Please enter a valid number for room size.")
		return nil
	}
	if MIN_ROOM_SIZE > roomSize || roomSize > MAX_ROOM_SIZE {
		client.WriteText(fmt.Sprintf("Please enter a room size between %d and %d.", MIN_ROOM_SIZE, MAX_ROOM_SIZE))
		return nil
	}
	room := models.CreateRoom(roomName, roomSize)
	rooms.Lock()
	defer rooms.Unlock()
	if len(rooms.roomMap) == MAX_ROOMS {
		client.WriteText("Max number of rooms reached. Please join an existing room or try again later.")
		return nil
	}
	rooms.roomMap[room.GetNumber()] = room
	log.Printf("(%s) %s created room \"%s\" with size %d", client.GetConn().RemoteAddr().String(), client.GetDisplayName(), roomName, roomSize)
	return room
}

// Runs a room in new goroutine, removes room when done
func runRoom(room *models.Room) {
	defer func() {
		rooms.Lock()
		defer rooms.Unlock()
		delete(rooms.roomMap, room.GetNumber())
	}()
	room.Run()
}

// print tabulated list of rooms to client
func listRooms(client *models.Client, rooms map[int]*models.Room) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Room Number\tRoom Name\tUsers")
	for _, room := range rooms {
		fmt.Fprintf(w, "%d\t%s\t%d/%d\n", room.GetNumber(), room.GetName(), room.GetHeadCount(), room.GetCapacity())
	}
	client.WriteText(buf.String())
}
