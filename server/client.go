package main

import "github.com/gorilla/websocket"

type Client struct {
	displayName string
	conn        *websocket.Conn
}

func (c *Client) writeText(msg string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func createClient(displayName string, conn *websocket.Conn) *Client {
	return &Client{
		displayName,
		conn,
	}
}
