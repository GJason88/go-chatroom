package utils

import (
	"github.com/gorilla/websocket"
)

var messageTypes = map[int]string{
	websocket.TextMessage:   "TextMessage",
	websocket.BinaryMessage: "BinaryMessage",
	websocket.CloseMessage:  "CloseMessage",
	websocket.PingMessage:   "PingMessage",
	websocket.PongMessage:   "PongMessage",
}

func FormatMessageType(msgType int) string {
	return messageTypes[msgType]
}
