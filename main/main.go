package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

var rdb *redis.Client

func main() {

	opt, err := redis.ParseURL("redis:127.0.0.1:6379") //ParseURL parses an URL into Options that can be used to connect to Redis by sockets
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(opt) // create an instance of a client
	http.HandleFunc("/join-chat", joinChat)
	go handleMessages() //set up a goroutine that decides what to do whenever a user sends a message
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
