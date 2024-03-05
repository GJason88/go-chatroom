package models

import (
	"chatroom/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	displayName string
	conn        *websocket.Conn
}

func (c *Client) ReadText() (string, error) {
	_, msgBytes, err := c.conn.ReadMessage()
	if err != nil {
		utils.LogReadErrors(err)
		return "", err
	}
	return string(msgBytes), nil
}

func (c *Client) WriteText(msg string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (c *Client) Help() {
	c.WriteText("TODO: Help action")
}

func (c *Client) GetDisplayName() string {
	return c.displayName
}

func (c *Client) GetRemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Client) Disconnect() (addr, displayName string) {
	addr = c.conn.RemoteAddr().String()
	displayName = c.displayName
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.conn.Close()
	return
}

func CreateClient(displayName string, conn *websocket.Conn) *Client {
	return &Client{
		displayName,
		conn,
	}
}
