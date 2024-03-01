package main

import (
	"bufio"
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
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if err == websocket.ErrCloseSent {
				log.Println("You have disconnected from the server")
			} else {
				log.Println("read error:", err)
			}
			break
		}
		log.Println(string(msgBytes))
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

	fmt.Print("Enter a display name: ")
	var displayName string
	_, err := fmt.Scanln(&displayName)
	if err != nil {
		log.Println("error:", err)
		return
	}

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/join", RawQuery: "displayName=" + displayName}
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
