package localServer;

import "github.com/gorilla/websocket"

/*
	CLIENT MANAGER STRUCT:
	clients: of type map where the keys are pointers to client structs and the values are boolean
	broadcast: of type chan (channel) used to broadcast byte data
	register: of type chan (channel) used to register the client, it sends the client's pointer using the channel
	unregister: of type chan (channel) used to unregister the client, it sends the client's pointer using the channel

*/
type ClientManager struct {
	clients 	map[*Client]bool
	broadcast 	chan []byte
	register 	chan *Client
	unregister 	chan *Client
}

/*
	CLIENT STRUCT:
	id: of type string, it stores the client's ID
	socket: of type pointer to a websocket connection, it holds the reference to the connection of the client
	send: of type chan (channel) of byte slices for sending data to the client
*/
type Client struct {
	id 		string
	socket 	*websocket.Conn
	send 	chan []byte
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

var manager = ClientManager {
	broadcast: make(chan []byte),
	register: make(chan *Client),
	unregister: make(chan *Client),
	clients: make(map[*Client]bool),
}