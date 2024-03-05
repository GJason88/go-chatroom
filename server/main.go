package main

import (
	"chatroom/server/models"
	"log"
	"net/http"
	"sync"
)

type Rooms struct {
	sync.Mutex
	roomMap map[int]*models.Room
}

type Clients struct {
	sync.Mutex
	clientMap map[string]*models.Client
}

var clients = Clients{clientMap: make(map[string]*models.Client)} // {remote addr: client} all clients in server
var rooms = Rooms{roomMap: make(map[int]*models.Room)}            // {int: room} all rooms in server

func main() {
	http.HandleFunc("/", connectClient) // connect client to lobby
	log.Fatal(http.ListenAndServe(":8080", nil))
}
