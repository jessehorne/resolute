package rhandlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

type ResetRoomKeysReqData struct {
	RoomID string `json:"room_id"`
}

type ResetRoomKeysReq struct {
	Cmd  string               `json:"cmd"`
	Data ResetRoomKeysReqData `json:"data"`
}

type ResetRoomKeysResData struct {
	RoomID string `json:"room_id"`
}

type ResetRoomKeysRes struct {
	Cmd  string               `json:"cmd"`
	Data ResetRoomKeysResData `json:"data"`
}

func ResetRoomKeys(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
	var r ResetRoomKeysReq
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	// check if room exists
	room, ok := s.Rooms[r.Data.RoomID]
	if !ok {
		c.WriteJSON(ResponseError{
			Cmd: "reset-room-keys-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "no room",
			},
		})
		return nil
	}

	// check perms
	if room.OwnerID != userID {
		c.WriteJSON(ResponseError{
			Cmd: "join-room-forever-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "unauthorized",
			},
		})
		return nil
	}

	// if it exists, reset keys to default (no more onetime or forever keys)
	room.ForeverJoinKey = ""
	room.OneTimeJoinKeys = []string{}

	// send response
	c.WriteJSON(ResetRoomKeysRes{
		Cmd: "reset-room-keys",
		Data: ResetRoomKeysResData{
			RoomID: r.Data.RoomID,
		},
	})
	return nil
}
