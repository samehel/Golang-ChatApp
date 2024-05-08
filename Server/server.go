package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

/*
CLIENT MANAGER STRUCT:
clients: of type map where the keys are pointers to client structs and the values are boolean
broadcast: of type chan (channel) used to broadcast byte data
register: of type chan (channel) used to register the client, it sends the client's pointer using the channel
unregister: of type chan (channel) used to unregister the client, it sends the client's pointer using the channel
*/
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

/*
CLIENT STRUCT:
id: of type string, it stores the client's ID
socket: of type pointer to a websocket connection, it holds the reference to the connection of the client
send: of type chan (channel) of byte slices for sending data to the client
*/
type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

/*
MESSAGE STRUCT:
sender: of type string which stores the sender of the message and has a struct tag to hold metadata about the sender in JSON format
recipient: of type string which stores the recipient of the message and has a struct tag to hold metadata about the recipient in JSON format
content: of type string which stores the content of the message and has a struct tag to hold metadata about the content of the message in JSON format
*/
type Message struct {
	sender    string `json:"sender,omitempty"`
	recipient string `json:"recipient,omitempty"`
	content   string `json:"content,omitempty"`
}

// Our global manager for our app
var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

/*
MAIN SERVER GOROUTINE:
Our server goroutine to be used for registering and
unregistering the client and handling the broadcast
and will run indefinitely until we tell it to stop.

whenever client is attempting to be registered, they
will be added to a map of available clients then a
message is sent to all clients to notify them that
someone new has connected.

whenever client is attempting to be unregistered, they
will be removed to a map of available clients then a
message is sent to all clients to notify them that
someone new has disconnected.

if the broadcast contains data it means someone is trying
to send and receive a message. So we will loop over all clients
and send their message and display it to them unless if
the channel is clogged then we will remove that client instead
*/
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{content: "Someone successfully connected"})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{content: "Someone successfully disconnected"})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

/*
A Global function used to send the client's message.
It loops over all the clients and checks whether the
the message is being sent by that client and if it is
then it will send the message
*/
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

/*
READ GOROUTINE:
This function is made for the purpose of reading data
from the socket and adding it to our broadcast for further
processing. In the case an error has occured, we assume the client
has disconnected and we unregister that client.
*/
func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{sender: c.id, content: string(message)})
		manager.broadcast <- jsonMessage
	}
}

/*
WRITE GOROUTINE:
This function is made for the purpose of writing data
from the socket. If for some reason it has issues reading
the message, we will send a message notifying the client they
have been disconnected
*/
func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

/*
Our main function to run all of our goroutines
*/
func main() {
	fmt.Println("Starting server...")
	go manager.start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}

/*
Our endpoint function
*/
func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	client_id, err := uuid.NewV7()
	if err != nil {
		fmt.Println("Error generating client id", err)
		return
	}

	client := &Client{id: client_id.String(), socket: conn, send: make(chan []byte)}

	manager.register <- client

	go client.read()
	go client.write()
}
