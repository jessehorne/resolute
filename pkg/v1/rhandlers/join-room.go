package rhandlers

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
	"github.com/jessehorne/resolute/pkg/v1/util"
)

type JoinRoomOneTimeReqData struct {
	RoomID     string `json:"room_id"`
	OneTimeKey string `json:"one_time_key"`
	Username   string `json:"username"`
	PublicKey  string `json:"public_key"`
}

type JoinRoomOneTimeReq struct {
	Cmd  string                 `json:"cmd"`
	Data JoinRoomOneTimeReqData `json:"data"`
}

type JoinRoomOneTimeResData struct {
	RoomID    string `json:"room_id"`
	RoomName  string `json:"room_name"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type JoinRoomOneTimeRes struct {
	Cmd   string                 `json:"cmd"`
	Data  JoinRoomOneTimeResData `json:"data"`
	Users []rstructs.JoinedUser  `json:"users"`
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
			key, err := util.ParsePublicKey(r.Data.PublicKey)
			if err != nil {
				// TODO?
				break
			}

			// tell current room users you joined
			for _, u := range room.Users {
				if u.UserID == userID {
					continue
				}

				u.Conn.WriteJSON(JoinRoomOneTimeRes{
					Cmd: "user-joined-onetime",
					Data: JoinRoomOneTimeResData{
						RoomID:    room.ID,
						RoomName:  room.Name,
						UserID:    userID,
						Username:  r.Data.Username,
						PublicKey: r.Data.PublicKey,
					},
				})
			}

			// tell yourself your id
			c.WriteJSON(YouJoinedRoomRes{
				Cmd: "joined",
				Data: YouJoinedRoomResData{
					UserID:   userID,
					RoomID:   room.ID,
					RoomName: room.Name,
					Users:    room.GetUsers(),
				},
			})

			newUser := &rstructs.User{
				UserID:    userID,
				Conn:      c,
				Username:  r.Data.Username,
				PublicKey: key,
			}
			room.AddUser(newUser)

			// delete one time key
			uptil := room.OneTimeJoinKeys[:i]
			until := room.OneTimeJoinKeys[i+1:]
			room.OneTimeJoinKeys = append(uptil, until...)

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
	PublicKey  string `json:"public_key"`
}

type JoinRoomForeverReq struct {
	Cmd  string                 `json:"cmd"`
	Data JoinRoomForeverReqData `json:"data"`
}

type JoinRoomForeverResData struct {
	RoomID    string                `json:"room_id"`
	RoomName  string                `json:"room_name"`
	UserID    string                `json:"user_id"`
	Username  string                `json:"username"`
	PublicKey string                `json:"public_key"`
	Users     []rstructs.JoinedUser `json:"users"`
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

	// add user to room and send response to all users in room
	key, err := util.ParsePublicKey(r.Data.PublicKey)
	if err != nil {
		// TODO
		return nil
	}

	// tell current room users you joined
	for _, u := range room.Users {
		if u.UserID == userID {
			continue
		}

		u.Conn.WriteJSON(JoinRoomForeverRes{
			Cmd: "user-joined-forever",
			Data: JoinRoomForeverResData{
				RoomID:    room.ID,
				RoomName:  room.Name,
				UserID:    userID,
				Username:  r.Data.Username,
				PublicKey: r.Data.PublicKey,
			},
		})
	}

	// tell yourself your id
	c.WriteJSON(YouJoinedRoomRes{
		Cmd: "joined",
		Data: YouJoinedRoomResData{
			UserID:   userID,
			RoomID:   room.ID,
			RoomName: room.Name,
			Users:    room.GetUsers(),
		},
	})

	newUser := &rstructs.User{
		UserID:    userID,
		Conn:      c,
		Username:  r.Data.Username,
		PublicKey: key,
	}
	room.AddUser(newUser)

	return nil
}

type YouJoinedRoomResData struct {
	UserID   string                `json:"user_id"`
	RoomID   string                `json:"room_id"`
	RoomName string                `json:"room_name"`
	Users    []rstructs.JoinedUser `json:"users"`
}

type YouJoinedRoomRes struct {
	Cmd  string               `json:"cmd"`
	Data YouJoinedRoomResData `json:"data"`
}
