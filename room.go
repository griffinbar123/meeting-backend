package main

import (
	"log"
)

var RoomArray = []Room{}

func NewRoom() Room {
	var roomNumber int
GetRoomNumberLoop:
	for {
		roomNumber = RangeIn(1000, 9999)
		for _, room := range RoomArray { //makes sure room id is unique
			if room.RoomId == roomNumber {
				continue GetRoomNumberLoop
			}
		}
		break GetRoomNumberLoop
	}
	return Room{
		RoomId:        roomNumber,
		RoomSize:      0,
		ActiveClients: make([]ActiveClient, 0),
	}
}

type Room struct { // stores data about a room and its active clients
	RoomId        int            `json:"room-id"`
	RoomSize      int            `json:"room-size"`
	ActiveClients []ActiveClient `json:"active-clients"`
}

func (room *Room) UpdateRoomForAllActiveClients() {
	/*
		Goes through all the active clients and sends them the current state of the room
	*/
	log.Println("Updating all active clients")
	for _, p := range room.ActiveClients {
		log.Println("Updating an active client: ", p.Name)
		log.Println(p.Times)
		p.SendRoomState(room)
	}
}

type RoomSender struct {
	Type    string `json:"type"`
	Payload Room   `json:"payload"`
}

func FindRoomFromId(id int) (int, *Room, bool) {
	for i, room := range RoomArray {
		if room.RoomId == id {
			return i, &room, true
		}
	}
	return 0, &Room{}, false
}

type RoomNumberSender struct {
	Type    string     `json:"type"`
	Payload RoomNumber `json:"payload"`
}

type RoomNumber struct {
	RoomNo int `json:"room-number"`
}
