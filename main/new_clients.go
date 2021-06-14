package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct { // data structure that captures user’s name and their message
	Username string `json:"username"`
	Text     string `json:"text"`
}

var clients = make(map[*websocket.Conn]bool) //list of all the currently active clients (or open WebSockets)
var broadcaster = make(chan Message)         // single channel that is responsible for sending and receiving Message data structure
var upgrader = websocket.Upgrader{           //necessary to “upgrade” Gorilla’s incoming requests into a WebSocket connection.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func joinChat(w http.ResponseWriter, r *http.Request) {
	//Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true //append a new client to clients map

	// if it's zero, no messages were ever sent/saved
	if rdb.Exists(context.Background(), "messages").Val() != 0 {
		sendPreviousMessages(ws)
	}

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		// send new message to the channel
		broadcaster <- msg
	}
}

func sendPreviousMessages(ws *websocket.Conn) {

	messages, err := rdb.LRange(context.Background(), "messages", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	// send previous messages
	for _, message := range messages {
		var msg Message
		json.Unmarshal([]byte(message), &msg)
		messageClient(ws, msg)
	}
}
