package main

import (
	"github.com/gorilla/websocket"
)

type ActiveClient struct { // an active client is a client with extra game specific data
	Client
	Conn     *websocket.Conn `json:"-"`
	Times    []Time          `json:"times"`
	ColorPos int             `json:"color-pos"`
}

func (activeClient *ActiveClient) SendRoomState(room *Room) {
	/*
		sends the room state to an active client
	*/
	// log.Println("Updating an active clientsss: ", activeClient.Times[0].StartDateTime)
	activeClient.Conn.WriteJSON(RoomSender{
		Type:    "room-state",
		Payload: *room,
	})
}

type Time struct {
	StartDateTime string `json:"start-date-time"`
	EndDateTime   string `json:"end-date-time"`
}
