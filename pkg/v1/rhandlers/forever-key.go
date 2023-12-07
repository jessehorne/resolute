package rhandlers

import (
	"encoding/json"

	"github.com/dchest/uniuri"
	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

type GetRoomForeverKeyRequest struct {
	Cmd  string    `json:"cmd"`
	Data RoomIDReq `json:"data"`
}

type RoomForeverKey struct {
	ForeverJoinKey string `json:"forever_join_key"`
}

type GetRoomForeverKeyResponse struct {
	Cmd  string                        `json:"cmd"`
	Data GetRoomForeverKeyResponseData `json:"data"`
}

type GetRoomForeverKeyResponseData struct {
	RoomID         string `json:"room_id"`
	ForeverJoinKey string `json:"forever_join_key"`
}

func GetRoomForeverKey(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
	var r GetRoomForeverKeyRequest
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	room, ok := s.Rooms[r.Data.RoomID]
	if !ok {
		c.WriteJSON(ResponseError{
			Cmd: "room-key-forever-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "no room",
			},
		})
		return nil
	}

	// check if user owns the room
	if s.Rooms[r.Data.RoomID].OwnerID != userID {
		c.WriteJSON(ResponseError{
			Cmd: "room-key-forever-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "unauthorized",
			},
		})
		return nil
	}

	if room.ForeverJoinKey == "" {
		room.ForeverJoinKey = uniuri.NewLen(32)
	}

	c.WriteJSON(GetRoomForeverKeyResponse{
		Cmd: "room-key-forever",
		Data: GetRoomForeverKeyResponseData{
			RoomID:         r.Data.RoomID,
			ForeverJoinKey: room.ForeverJoinKey,
		},
	})

	return nil
}
