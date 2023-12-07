package client

import (
	"errors"

	"github.com/jessehorne/resolute/pkg/v1/rhandlers"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
	"github.com/jessehorne/resolute/pkg/v1/util"
)

type CRoom struct {
	ClientRoomID string
	Client       *Client

	IsOwner     bool
	ReqName     string
	ReqUsername string

	Room         *rstructs.Room
	MessageQueue []string

	callbacks map[string]interface{}
}

func (r *CRoom) UpdateUsers(users []rstructs.JoinedUser) {
	if r.Room.Users == nil {
		r.Room.Users = map[string]*rstructs.User{}
	}

	for _, u := range users {
		if u.UserID == r.Client.UserID {
			continue
		}

		pubKey, err := util.ParsePublicKey(u.PublicKeyString)
		if err != nil {
			continue
		}

		existingUser, ok := r.Room.Users[u.UserID]
		if !ok {
			existingUser = &rstructs.User{}
			r.Room.Users[u.UserID] = existingUser
		}
		existingUser.UserID = u.UserID
		existingUser.Username = u.Username
		existingUser.PublicKeyString = u.PublicKeyString
		existingUser.PublicKey = pubKey
	}
}

func (r *CRoom) AddJoinedUser(userID, username, pubKey string) {
	if r.Room.Users == nil {
		r.Room.Users = map[string]*rstructs.User{}
	}

	pk, err := util.ParsePublicKey(pubKey)
	if err != nil {
		return
	}

	existingUser, ok := r.Room.Users[userID]
	if !ok {
		existingUser = &rstructs.User{}
		r.Room.Users[userID] = existingUser
	}
	existingUser.UserID = userID
	existingUser.Username = username
	existingUser.PublicKeyString = pubKey
	existingUser.PublicKey = pk
}

func (r *CRoom) On(name string, cb interface{}) {
	// TODO: add checks to make sure the functions being provided are acceptable
	// example: key-onetime requires a function like func(string, string) since it returns (roomID, key)
	r.callbacks[name] = cb
}

func (r *CRoom) call(name string, data interface{}) {
	cb, ok := r.callbacks[name]
	if ok {
		if name == "created" {
			cb.(func(string))(data.(string))
		} else if name == "key-onetime" {
			d := data.(map[string]string)
			cb.(func(string, string))(d["room_id"], d["key"])
		} else if name == "key-forever" {
			d := data.(map[string]string)
			cb.(func(string, string))(d["room_id"], d["key"])
		} else if name == "send-message" {
			d := data.(map[string]string)
			cb.(func(string, string, string, string))(d["room_id"], d["user_id"], d["username"], d["content"])
		} else if name == "joined" {
			d := data.(map[string]string)
			cb.(func(string, string))(d["room_id"], d["room_name"])
		} else if name == "user-joined" {
			d := data.(map[string]string)
			cb.(func(string, string, string, string, string))(d["room_id"], d["room_name"], d["user_id"], d["username"], d["key_type"])
		}
	}
}

func (r *CRoom) GetKey(t string) error {
	if r.Room == nil {
		return errors.New("room not ready")
	}

	if t == "onetime" {
		req := rhandlers.GetRoomOneTimeKeyRequest{
			Cmd: "room-key-onetime",
			Data: rhandlers.RoomIDReq{
				RoomID: r.Room.ID,
			},
		}

		err := r.Client.Conn.WriteJSON(req)
		if err != nil {
			// TODO
		}
	} else if t == "forever" {
		req := rhandlers.GetRoomForeverKeyRequest{
			Cmd: "room-key-forever",
			Data: rhandlers.RoomIDReq{
				RoomID: r.Room.ID,
			},
		}

		err := r.Client.Conn.WriteJSON(req)
		if err != nil {
			// TODO
		}
	}

	return nil
}

func (r *CRoom) SendMessage(content string) error {
	if r.Room == nil {
		return errors.New("room not set")
	}

	// send message to all users in room encrypted with their public key
	for _, u := range r.Room.Users {
		if u.UserID == r.Client.UserID {
			continue
		}

		if u.PublicKey == nil {
			continue
		}

		s, err := util.EncryptMessage(u.PublicKey, content)
		if err != nil {
			continue
		}

		req := rhandlers.SendMessageReq{
			Cmd: "send-message",
			Data: rhandlers.SendMessageReqData{
				RoomID:   r.Room.ID,
				ToUserID: u.UserID,
				Content:  s,
			},
		}

		err = r.Client.Conn.WriteJSON(req)
		if err != nil {
			// TODO
		}
	}

	return nil
}
