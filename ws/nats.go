package ws

import (
	"encoding/json"
	"log"
	"ws_realtime_app/lib"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func InitNATS(natsURL string) *nats.Conn {
	/* Initializes NATS connection and starts message listener */
	var err error
	if nc, err = nats.Connect(natsURL); err != nil {
		log.Println("Error connecting to NATS:", err)
	}

	InitMsgListener()

	return nc
}

func InitMsgListener() {
	/* This function will receive messages fromNATS and dispatch to all sockets in rooms */
	nc.Subscribe("message", func(m *nats.Msg) {
		// TODO: Process message and dispatch to connected sockets and rooms
		log.Printf("Received NATS msg: %s\n", string(m.Data))

		// Unmarshals received JSON into model for further processing
		// TODO: In case further message verification and integrity needs to be checked, do it here
		var msg lib.MessageModel
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			panic(err)
		}

		// Now we will start dispatching it to all sockets
		sock_list := GetSocketsForRoom(msg.RoomID)

		// Looping over all sockets in the room and forwarding the message to them
		for _, s := range sock_list {
			if err := s.WsConn.WriteJSON(msg); err != nil {
				// Closing connection in case of write error, client-side should reconnect
				defer s.WsConn.Close()
			}
		}
	})
}

func SendMsgNATS(content any) {
	/* Dispatch new message to NATS cluster */

	payload, err := json.Marshal(content)
	if err == nil {
		nc.Publish("message", payload)
	}
}
