package ws

import (
	"golang.org/x/net/websocket"
	"fmt"
	"strconv"
)

var clients = make([]*websocket.Conn, 0)

// new websocket connection handler
func EchoHandler(ws *websocket.Conn) {
	var err error

	clients = append(clients, ws)
	fmt.Println("Total clients " + strconv.Itoa(len(clients)))

	for {

		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive", err.Error())
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}
