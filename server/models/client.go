package models

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	DisplayName string
	Conn        *websocket.Conn
}

func (c *Client) WriteText(msg string) {
	c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *Client) Help() {
	c.WriteText("TODO: Help action")
}

func CreateClient(displayName string, conn *websocket.Conn) *Client {
	return &Client{
		displayName,
		conn,
	}
}
