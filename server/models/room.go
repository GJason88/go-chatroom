package models

import (
	"fmt"
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

func (r *Room) AddClient(client *Client) {
	defer func() { go r.listen(client) }()
	r.clients.Lock()
	defer r.clients.Unlock()
	if r.GetHeadCount() == 0 {
		client.WriteText("Room no longer exists.")
		return
	}
	r.clients.clientMap[client.GetRemoteAddr()] = client
	client.WriteText(fmt.Sprintf("You have connected to \"%s\". Type \"/leave\" to leave.", r.name))
	r.messages <- fmt.Sprintf("%s has joined the room", client.GetDisplayName())
}

func (r *Room) listen(client *Client) {
	defer r.removeClient(client)
	if _, ok := r.clients.clientMap[client.GetRemoteAddr()]; !ok {
		return
	}
	for {
		msg, err := client.ReadText()
		if err != nil || msg == "/leave" {
			break
		}
		r.messages <- fmt.Sprintf("%s: %s", client.GetDisplayName(), msg)
	}
}

// TODO: solution to relisten to client in lobby
func (r *Room) removeClient(client *Client) {
	delete(r.clients.clientMap, client.GetRemoteAddr())
	r.messages <- fmt.Sprintf("%s has left the room", client.GetDisplayName())
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
