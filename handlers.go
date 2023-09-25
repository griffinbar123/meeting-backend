package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func EstablishWS(w http.ResponseWriter, r *http.Request) {
	/*
	 upgrades the http connection to a websocket connection, then listens for messages and handles them
	*/
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		HandleMessage(ws, message)
	}
}
func HandleMessage(ws *websocket.Conn, message []byte) {
	/*
		 handles the different type od messages the client sends. messages are formatted:
		 {
			"type": "do-something",
			"payload": {"data": "data"}
		 }
	*/
	messageAsString := string(message)
	j := objx.MustFromJSON(messageAsString)
	payload := j.Exclude([]string{"type"})
	temp, err := payload.JSON()
	if err != nil {
		log.Printf("error: %v", err)
	}
	t := j.Get("type").Str()
	log.Println("message: " + t)

	switch t {
	case "active-client-metadata":
		var client = Client{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &client); err != nil {
			log.Printf("error: %v", err)
		}
		if activeClient, ok := client.CheckForClientMetadata(ws); ok {
			activeClient.JoinRoom()
		}
	case "join-room":
		var client = ClientInRoom{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &client); err != nil {
			log.Printf("error: %v", err)
		}
		client.Conn = ws
		client.JoinRoom()
	case "create-room":
		room := NewRoom()
		RoomArray = append(RoomArray, room)
		ws.WriteJSON(RoomNumberSender{
			Type: "send-room-number",
			Payload: RoomNumber{
				RoomNo: room.RoomId,
			},
		})
	case "disconnect":
		var client = Client{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &client); err != nil {
			log.Printf("error: %v", err)
		}

		client.DisconnectClientFromTheirRoom()
		if err := ws.Close(); err != nil {
			log.Printf("error: %v", err)
		}

	case "room-state":
		var activeClientRoomState = Room{}
		if json.Unmarshal([]byte(temp[11:len(temp)-1]), &activeClientRoomState); err != nil {
			log.Printf("error: %v", err)
		}
		log.Println("room-state: " + temp[11:len(temp)-1])

		if i, _, ok := FindRoomFromId(activeClientRoomState.RoomId); ok {
			for j, client := range activeClientRoomState.ActiveClients {
				RoomArray[i].ActiveClients[j].Times = client.Times
			}
			RoomArray[i].UpdateRoomForAllActiveClients()
		}

	default:
		log.Println("error: unknown message type")
	}

}
