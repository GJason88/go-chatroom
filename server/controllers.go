package main

import (
	"chatroom/server/models"
	"strconv"
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
	client.WriteText("TODO: list rooms")
}
