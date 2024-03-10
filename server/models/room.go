package models

import (
	"fmt"
	"log"
	"sync"
)

type Clients struct {
	sync.Mutex
	clientMap map[string]*Client
}

type Room struct {
	number   int
	name     string
	capacity int
	clients  Clients
	messages chan string
	done     chan bool
}

func (r *Room) writeToClients(msg string) {
	for _, c := range r.clients.clientMap {
		c.WriteText(msg)
	}
}

func (r *Room) AddClient(client *Client, isCreator bool) {
	defer r.removeClient(client)
	r.clients.Lock()
	if !isCreator && r.GetHeadCount() == 0 {
		r.clients.Unlock()
		client.WriteText("Room is unavailable.")
		return
	}
	if len(r.clients.clientMap) == r.capacity {
		r.clients.Unlock()
		client.WriteText("Room is full.")
		return
	}
	r.clients.clientMap[client.GetRemoteAddr()] = client
	r.clients.Unlock()
	log.Printf("(%s) %s has joined room %d (%s)", client.GetRemoteAddr(), client.GetDisplayName(), r.number, r.name)
	r.messages <- fmt.Sprintf("%s has joined the room", client.GetDisplayName())
	r.listen(client) // blocks
}

func (r *Room) listen(client *Client) {
	client.WriteText(fmt.Sprintf("You have connected to \"%s\". Type \"/leave\" to leave.", r.name))
	for {
		msg, err := client.ReadText()
		if err != nil || msg == "/leave" {
			break
		}
		r.messages <- fmt.Sprintf("%s: %s", client.GetDisplayName(), msg)
	}
}

func (r *Room) removeClient(client *Client) {
	if !r.HasClient(client) {
		return
	}
	r.clients.Lock()
	defer r.clients.Unlock()
	delete(r.clients.clientMap, client.GetRemoteAddr())
	log.Printf("(%s) %s has left room %d (%s)", client.GetRemoteAddr(), client.GetDisplayName(), r.number, r.name)
	r.messages <- fmt.Sprintf("%s has left the room", client.GetDisplayName())
	if len(r.clients.clientMap) == 0 {
		r.done <- true
	}
}

func (r *Room) HasClient(client *Client) bool {
	_, ok := r.clients.clientMap[client.GetRemoteAddr()]
	return ok
}

func (r *Room) GetCapacity() int {
	return r.capacity
}

func (r *Room) GetNumber() int {
	return r.number
}

func (r *Room) GetName() string {
	return r.name
}

func (r *Room) GetHeadCount() int {
	return len(r.clients.clientMap)
}

func (r *Room) Run() {
	for {
		select {
		case <-r.done:
			return
		case msg := <-r.messages:
			r.writeToClients(msg)
		}
	}
}

func CreateRoom(roomName string, size int) *Room {
	autoIncRoomId.Lock()
	defer autoIncRoomId.Unlock()
	autoIncRoomId.id++
	return &Room{
		autoIncRoomId.id,
		roomName,
		size,
		Clients{clientMap: make(map[string]*Client)},
		make(chan string),
		make(chan bool),
	}
}
