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

func HandleError(err error, msg string, shouldPanic bool) {
	if err != nil {
		if shouldPanic {
			panic(err)
		}
		log.Printf("%s: %v", msg, err)
	}
}
