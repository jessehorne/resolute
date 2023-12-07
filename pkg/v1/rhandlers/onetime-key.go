package rhandlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

type GetRoomOneTimeKeyRequest struct {
	Cmd  string    `json:"cmd"`
	Data RoomIDReq `json:"data"`
}

type RoomOneTimeKey struct {
	RoomID     string `json:"room_id"`
	OneTimeKey string `json:"one_time_key"`
}

type GetRoomOneTimeKeyResponse struct {
	Cmd  string         `json:"cmd"`
	Data RoomOneTimeKey `json:"data"`
}

func GetRoomOneTimeKey(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
	var r GetRoomOneTimeKeyRequest
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	// check if state has room by RoomID
	if !s.HasRoom(r.Data.RoomID) {
		c.WriteJSON(ResponseError{
			Cmd: "room-key-onetime-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "room doesn't exist",
			},
		})
		return nil
	}

	// check if user owns the room
	if s.Rooms[r.Data.RoomID].OwnerID != userID {
		c.WriteJSON(ResponseError{
			Cmd: "room-key-onetime-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "unauthorized",
			},
		})
	}

	key := s.CreateOneTimeRoomKey(r.Data.RoomID)

	c.WriteJSON(GetRoomOneTimeKeyResponse{
		Cmd: "room-key-onetime",
		Data: RoomOneTimeKey{
			RoomID:     r.Data.RoomID,
			OneTimeKey: key,
		},
	})

	return nil
}
