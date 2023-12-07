package rhandlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

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
