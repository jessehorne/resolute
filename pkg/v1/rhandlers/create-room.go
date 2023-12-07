package rhandlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
	"github.com/jessehorne/resolute/pkg/v1/util"
)

type CreateRoomRequestData struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	ClientRoomID    string `json:"client_room_id"`
	PublicKeyString string `json:"public_key_string"`
}

type CreateRoomRequest struct {
	Cmd  string                `json:"cmd"`
	Data CreateRoomRequestData `json:"data"`
}

type CreateRoomResponse struct {
	Cmd          string         `json:"cmd"`
	Room         *rstructs.Room `json:"room"`
	ClientRoomID string         `json:"client_room_id"`
}

func CreateRoomHandler(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
	var cr CreateRoomRequest
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return err
	}

	pubKey, err := util.ParsePublicKey(cr.Data.PublicKeyString)
	if err != nil {
		// TODO?
		return nil
	}

	newRoom := rstructs.NewRoom(cr.Data.Name, userID)
	newRoom.OwnerID = userID
	newRoom.AddUser(&rstructs.User{
		UserID:    userID,
		Username:  cr.Data.Username,
		Conn:      c,
		PublicKey: pubKey,
	})
	s.AddRoom(newRoom)

	c.WriteJSON(CreateRoomResponse{
		Cmd:          "create-room",
		Room:         newRoom,
		ClientRoomID: cr.Data.ClientRoomID,
	})

	return nil
}
