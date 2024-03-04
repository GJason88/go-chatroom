package models

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	displayName string
	conn        *websocket.Conn
}

func (c *Client) WriteText(msg string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *Client) Help() {
	c.WriteText("TODO: Help action")
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

func (c *Client) GetDisplayName() string {
	return c.displayName
}

func CreateClient(displayName string, conn *websocket.Conn) *Client {
	return &Client{
		displayName,
		conn,
	}
}
