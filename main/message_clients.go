package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
)

//sending any new messages to every connected client
func handleMessages() {
	for {
		// grab any next message from channel
		msg := <-broadcaster

		storeInRedis(msg)
		// send to every client currently connected
		for client := range clients {
			messageClient(client, msg)
		}
	}
}

func storeInRedis(msg Message) {
	json, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	// store messages in Redis
	err = rdb.RPush(context.Background(), "messages", json).Err()
	if err != nil {
		panic(err)
	}
}

func messageClient(client *websocket.Conn, msg Message) {
	err := client.WriteJSON(msg)
	//If there’s an issue, print a message, close the client, and remove it from the map
	if err != nil {
		//log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
}
