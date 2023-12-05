package client

import (
	"errors"

	"github.com/jessehorne/resolute/pkg/v1/rhandlers"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

type CRoom struct {
	Client *Client

	IsOwner     bool
	ReqName     string
	ReqUsername string

	Room         *rstructs.Room
	MessageQueue []string

	callbacks map[string]interface{}
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
			cb.(func())()
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

func (r *CRoom) SendMessage(content string) {
	if r.Room == nil {
		return
	}

	req := rhandlers.SendMessageReq{
		Cmd: "send-message",
		Data: rhandlers.SendMessageReqData{
			RoomID:  r.Room.ID,
			Content: content,
		},
	}

	err := r.Client.Conn.WriteJSON(req)
	if err != nil {
		// TODO
	}
}
