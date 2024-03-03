package utils

import (
	"log"

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

func LogReadErrors(err error) {
	if !(websocket.IsCloseError(err, 1000) || err == websocket.ErrCloseSent) {
		log.Println("read error:", err)
	}
}
