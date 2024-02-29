package main

import (
	"bufio"
	"chatroom/utils"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func handleReads(conn *websocket.Conn, closeFlag chan struct{}) {
	defer close(closeFlag)
	for {
		fmt.Print("> ")
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("message: %s, type: %s", msgBytes, utils.FormatMessageType(msgType))
	}
}

func handleWrites(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Print("> ")
			continue
		}
		err := conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("write error:", err)
		}
	}

}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/join", RawQuery: "displayName=testuser"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial error:", err)
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
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then waiting (with timeout) for the server to close the connection.
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
				return
			}
			// block until closeFlag or 3 second timeout
			select {
			case <-closeFlag:
			case <-time.After(3 * time.Second):
			}
			return
		}
	}
}
