# chat-project
Install redis server https://redis.io/download and run it on your machine `redis-server`
We will use Redis to store messages that clients send to a chat and when a new client joins chat we will read all stored messages and send them to a client.

To test it use postman or any other app that can work with WebSockets.
For client-server communication we will use WebSockets.

Join chat: `ws://localhost:80/join-chat`

Send message:
`{
"username": "User 1",
"text": "Hello world!"
}`