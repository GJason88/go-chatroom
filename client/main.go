package main

import (
	"bufio"
	"chatroom/utils"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func handleReads(conn *websocket.Conn, closeFlag chan struct{}) {
	defer close(closeFlag)
	for {
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			if err == websocket.ErrCloseSent {
				break
			}
		}
		log.Printf("message: %s, type: %s", msgBytes, utils.FormatMessageType(msgType))
	}
}

func handleWrites(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		err := conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/echo", nil)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer conn.Close()

	closeFlag := make(chan struct{})
	go handleReads(conn, closeFlag)
	go handleWrites(conn)

	for {
		select {
		case <-closeFlag:
			return
		case <-interrupt:
			log.Println("interrupt:")
			// Cleanly close the connection by sending a close message and then waiting (with timeout) for the server to close the connection.
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			select {
			case <-closeFlag:
			case <-time.After(3 * time.Second):
			}
		}
	}
}
