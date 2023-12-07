package rhandlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

type CreateRoomRequestData struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type CreateRoomRequest struct {
	Cmd  string                `json:"cmd"`
	Data CreateRoomRequestData `json:"data"`
}

type CreateRoomResponse struct {
	Cmd  string         `json:"cmd"`
	Room *rstructs.Room `json:"room"`
}

func CreateRoomHandler(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
	var cr CreateRoomRequest
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return err
	}

	newRoom := rstructs.NewRoom(cr.Data.Name, userID)
	newRoom.AddUser(&rstructs.User{
		UserID:   userID,
		Username: cr.Data.Username,
		Conn:     c,
	})
	s.AddRoom(newRoom)

	c.WriteJSON(CreateRoomResponse{
		Cmd:  "create-room",
		Room: newRoom,
	})

	return nil
}
