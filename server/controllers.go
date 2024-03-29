package main

import (
	"bytes"
	"chatroom/server/models"
	"fmt"
	"log"
	"strconv"
	"text/tabwriter"
)

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
	log.Printf("(%s) %s created room \"%s\" with capacity %d", client.GetRemoteAddr(), client.GetDisplayName(), roomName, roomSize)
	return room
}

// Runs a room, removes room when done
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
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Room Number\tRoom Name\tUsers")
	for _, room := range rooms {
		fmt.Fprintf(w, "%d\t%s\t%d/%d\n", room.GetNumber(), room.GetName(), room.GetHeadCount(), room.GetCapacity())
	}
	if err := w.Flush(); err != nil {
		client.WriteText("Failed to list rooms.")
		log.Println(err)
		return
	}
	client.WriteText(buf.String())
}

func helpCommand(client *models.Client) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "A simple CLI chatroom application that supports up to %d clients at a time. Supports up to %d rooms at a time, each with a capacity of %d to %d clients.\n\nCommands:\nrooms\tList all existing rooms.\njoin [room_number]\tJoin an existing room.\ncreate [room_name] [capacity]\tCreate a room with a name and capacity between 2 and 8.\nquit\tDisconnect from the chatroom application.\n", SERVER_CAPACITY, MAX_ROOMS, MIN_ROOM_SIZE, MAX_ROOM_SIZE)
	if err := w.Flush(); err != nil {
		client.WriteText("Command failed.")
		log.Println(err)
		return
	}
	client.WriteText(buf.String())
}
