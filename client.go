package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func (client *Client) CheckForClientMetadata(conn *websocket.Conn) (ClientInRoom, bool) {
	/*
		looks at the localhost data send by the client and checks if the active client data it contains
		exists as an active active client connection
	*/
	if activeClient, ok := client.FindClientInList(); ok {
		activeClient.Conn = conn
		return activeClient, true
	}
	return ClientInRoom{}, false
}

func (client *ClientInRoom) RemoveClientFromClientArray() {
	for i, p := range ActiveClientArray {
		if p.Uuid == client.Uuid {
			ActiveClientArray = remove(ActiveClientArray, i)
		}
	}
}

func (client *Client) FindClientInList() (ClientInRoom, bool) {
	/*
		checks if a client is in the list of active clients
	*/
	for _, activeClient := range ActiveClientArray {
		if activeClient.Uuid == client.Uuid {
			return activeClient, true
		}
	}
	return ClientInRoom{}, false
}

func (client *Client) DisconnectClientFromTheirRoom() {
	/*
		takes a client and removes them from the room in which they are apart of
	*/
	if clientInRoom, ok := client.FindClientInList(); ok {
		if _, room, ok := FindRoomFromId(clientInRoom.RoomId); ok {
			for i, p := range room.ActiveClients {
				if p.Uuid == clientInRoom.Uuid {
					log.Println("Removing client from room: ", p.Name)
					room.ActiveClients = removeActiveClient(room.ActiveClients, i)
					room.RoomSize -= 1
					room.UpdateRoomForAllActiveClients()
					break
				}
			}
		}
	}
}

type ClientInRoom struct { //struct to store a client and their associeted connection and room id
	Client
	Conn   *websocket.Conn `json:"-"`
	RoomId int             `json:"room-id"`
}

type ClientSender struct {
	Type    string `json:"type"`
	Payload Client `json:"payload"`
}

func (client *ClientInRoom) AddClientToActiveClientArray() {
	ActiveClientArray = append(ActiveClientArray, *client)
}

func (client *ClientInRoom) JoinRoom() {
	/*
		handles adding a client to a room (if they are not already in the room)
	*/
	if i, room, ok := FindRoomFromId(client.RoomId); ok {
		if len(room.ActiveClients) > 8 {
			log.Printf("warning: cannot add active clients - room is already at max capacity")
			return
		}
		for j, p := range room.ActiveClients {
			if p.Uuid == client.Uuid {
				room.ActiveClients[j].Conn = client.Conn
				client.HandleActiveClientUpdate(room)
				return
			}
		}
		client.AddActiveClientToRoom(room)
		RoomArray[i] = *room
		client.StoreConnection()
		client.HandleActiveClientUpdate(room)
		return
	}
}

func (client *ClientInRoom) AddActiveClientToRoom(room *Room) {
	/*
		handles adding a client to a room and making them an active client
	*/
	room.ActiveClients = append(room.ActiveClients, client.NewActiveClient(room))
	room.RoomSize += 1
}

func (client *ClientInRoom) NewActiveClient(room *Room) ActiveClient {
	return ActiveClient{
		Client:   client.Client,
		Conn:     client.Conn,
		Times:    []Time{},
		ColorPos: room.RoomSize,
	}
}

func (client *ClientInRoom) StoreConnection() {
	/*
		stores active client info on client side in localstorage
	*/
	client.RemoveClientFromClientArray()
	client.AddClientToActiveClientArray()
	client.Conn.WriteJSON(ClientSender{
		Type:    "active-client-metadata",
		Payload: client.Client,
	})
}

func (client *ClientInRoom) HandleActiveClientUpdate(room *Room) {
	/*
		when a active client updates, we send the current roomdata to all the active clients so everyone
		is looking at the same room
	*/
	room.UpdateRoomForAllActiveClients()
}
