package handlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/structs"
)

type SendMessageReqData struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

type SendMessageReq struct {
	Cmd  string             `json:"cmd"`
	Data SendMessageReqData `json:"data"`
}

type SendMessageResData struct {
	RoomID  string `json:"room_id"`
	From    string `json:"from"`
	Content string `json:"content"`
}

type SendMessageRes struct {
	Cmd  string             `json:"cmd"`
	Data SendMessageResData `json:"data"`
}

func SendMessage(s *structs.State, userID string, c *websocket.Conn, data []byte) error {
	var r SendMessageReq
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	// check if room exists
	room, ok := s.Rooms[r.Data.RoomID]
	if !ok {
		c.WriteJSON(ResponseError{
			Cmd: "send-message-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "no room",
			},
		})
		return nil
	}

	// check if user is in room
	if !room.HasUser(userID) {
		c.WriteJSON(ResponseError{
			Cmd: "send-message-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "unauthorized",
			},
		})
		return nil
	}

	room.BroadcastMessage(userID, r.Data.Content)

	return nil
}
