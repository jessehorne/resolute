package rhandlers

import (
	"encoding/json"

	"github.com/dchest/uniuri"
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

type RoomIDReq struct {
	RoomID string `json:"room_id"`
}

type GetRoomOneTimeKeyRequest struct {
	Cmd  string    `json:"cmd"`
	Data RoomIDReq `json:"data"`
}

type ResponseError struct {
	Cmd  string `json:"cmd"`
	Data map[string]string
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

type JoinRoomOneTimeData struct {
	RoomID     string `json:"room_id"`
	OneTimeKey string `json:"one_time_key"`
	Username   string `json:"username"`
}

type JoinRoomOneTimeReq struct {
	Cmd  string              `json:"cmd"`
	Data JoinRoomOneTimeData `json:"data"`
}

type JoinRoomOneTimeResData struct {
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
}

type JoinRoomOneTimeRes struct {
	Cmd  string                 `json:"cmd"`
	Data JoinRoomOneTimeResData `json:"data"`
}

func JoinRoomOneTime(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
	var r JoinRoomOneTimeReq
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	// check that room exists
	room, ok := s.Rooms[r.Data.RoomID]
	if !ok {
		c.WriteJSON(ResponseError{
			Cmd: "join-room-onetime-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "no room",
			},
		})
		return nil
	}

	// check that one time key exists
	for i, k := range room.OneTimeJoinKeys {
		if k == r.Data.OneTimeKey {
			// add user to room
			room.AddUser(&rstructs.User{
				UserID:   userID,
				Username: r.Data.Username,
				Conn:     c,
			})

			// delete one time key
			uptil := room.OneTimeJoinKeys[:i]
			until := room.OneTimeJoinKeys[i+1:]
			room.OneTimeJoinKeys = append(uptil, until...)

			// send response
			c.WriteJSON(JoinRoomOneTimeRes{
				Cmd: "join-room-onetime",
				Data: JoinRoomOneTimeResData{
					RoomID:   r.Data.RoomID,
					RoomName: room.Name,
				},
			})

			return nil
		}
	}

	// no key was found
	c.WriteJSON(ResponseError{
		Cmd: "join-room-onetime-error",
		Data: map[string]string{
			"room_id": r.Data.RoomID,
			"msg":     "no onetime key",
		},
	})
	return nil
}

type JoinRoomForeverReqData struct {
	RoomID     string `json:"room_id"`
	ForeverKey string `json:"forever_key"`
	Username   string `json:"username"`
}

type JoinRoomForeverReq struct {
	Cmd  string                 `json:"cmd"`
	Data JoinRoomForeverReqData `json:"data"`
}

type JoinRoomForeverResData struct {
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
}

type JoinRoomForeverRes struct {
	Cmd  string                 `json:"cmd"`
	Data JoinRoomForeverResData `json:"data"`
}

func JoinRoomForever(s *rstructs.State, userID string, c *websocket.Conn, data []byte) error {
	var r JoinRoomForeverReq
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	// check that room exists
	room, ok := s.Rooms[r.Data.RoomID]
	if !ok {
		c.WriteJSON(ResponseError{
			Cmd: "join-room-forever-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "no room",
			},
		})
		return nil
	}

	// check if room forever key is "" and automatically deny
	if room.ForeverJoinKey == "" {
		c.WriteJSON(ResponseError{
			Cmd: "join-room-forever-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "unset forever key",
			},
		})
		return nil
	}

	// check that forever key is valid
	if room.ForeverJoinKey != r.Data.ForeverKey {
		c.WriteJSON(ResponseError{
			Cmd: "join-room-forever-error",
			Data: map[string]string{
				"room_id": r.Data.RoomID,
				"msg":     "no forever key",
			},
		})
		return nil
	}

	// add user to room and send response
	room.AddUser(&rstructs.User{
		UserID:   userID,
		Username: r.Data.Username,
		Conn:     c,
	})

	c.WriteJSON(JoinRoomForeverRes{
		Cmd: "join-room-forever",
		Data: JoinRoomForeverResData{
			RoomID:   room.ID,
			RoomName: room.Name,
		},
	})

	return nil
}

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
