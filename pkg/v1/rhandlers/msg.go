package rhandlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

type SendMessageReqData struct {
	RoomID   string `json:"room_id"`
	Content  string `json:"content"`
	ToUserID string `json:"to_user_id"`
}

type SendMessageReq struct {
	Cmd  string             `json:"cmd"`
	Data SendMessageReqData `json:"data"`
}

type SendMessageResData struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

type SendMessageRes struct {
	Cmd  string             `json:"cmd"`
	Data SendMessageResData `json:"data"`
}

func SendMessage(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
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

	room.BroadcastMessage(r.Data.ToUserID, r.Data.Content)

	return nil
}
