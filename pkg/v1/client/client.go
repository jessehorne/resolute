package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/pkg/v1/resolute"
	"github.com/jessehorne/resolute/pkg/v1/rhandlers"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

const (
	RoomKeyTypeOneTime = iota
	RoomKeyTypeForever
)

type Client struct {
	Host string
	Path string

	UserID string

	// when user creates room, we add CRoom to RoomQueue. When Rooms is filled,
	// we'll update the CRoom's reference to Room
	Rooms     map[string]*CRoom // rooms the user is in. key is the room id
	RoomQueue []*CRoom

	Conn *websocket.Conn
}

func NewClient(path, host string) *Client {
	u := url.URL{Scheme: "ws", Host: host, Path: path}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	newClient := &Client{
		Host:      host,
		Path:      path,
		Rooms:     map[string]*CRoom{},
		RoomQueue: []*CRoom{},
		Conn:      c,
	}

	return newClient
}

func (c *Client) Listen() {
	stopRunning := make(chan os.Signal, 1)
	signal.Notify(stopRunning, os.Interrupt)
	signal.Notify(stopRunning, syscall.SIGTERM)

	defer c.Conn.Close()

	go func() {
		for {
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				// TODO
				return
			}

			cmd, err := resolute.NewCommandFromJSON(message)
			if err != nil {
				// TODO
			}

			if cmd.Cmd == "create-room" {
				var r rhandlers.CreateRoomResponse
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				if len(c.RoomQueue) == 0 {
					continue
				}

				// get first room in queue
				first := c.RoomQueue[0]
				first.Room = &rstructs.Room{
					ID:              r.Room.ID,
					OwnerID:         r.Room.OwnerID,
					Name:            r.Room.Name,
					AdminKey:        "",
					OneTimeJoinKeys: []string{},
					ForeverJoinKey:  "",
					Users:           map[string]*rstructs.User{},
				}

				c.Rooms[r.Room.ID] = first
				c.RoomQueue = c.RoomQueue[1:]

				first.call("created", nil)
			} else if cmd.Cmd == "room-key-onetime" {
				var r rhandlers.GetRoomOneTimeKeyResponse
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				room, ok := c.Rooms[r.Data.RoomID]
				if ok {
					room.call("key-onetime", map[string]string{
						"room_id": r.Data.RoomID,
						"key":     r.Data.OneTimeKey,
					})
				}
			} else if cmd.Cmd == "room-key-forever" {
				var r rhandlers.GetRoomForeverKeyResponse
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				room, ok := c.Rooms[r.Data.RoomID]
				if ok {
					room.call("key-forever", map[string]string{
						"room_id": r.Data.RoomID,
						"key":     r.Data.ForeverJoinKey,
					})
				}
			} else if cmd.Cmd == "send-message" {
				var r rhandlers.SendMessageRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				room, ok := c.Rooms[r.Data.RoomID]
				if ok {
					room.call("send-message", map[string]string{
						"room_id":  r.Data.RoomID,
						"user_id":  r.Data.UserID,
						"username": r.Data.Username,
						"content":  r.Data.Content,
					})
				}
			} else if cmd.Cmd == "join-room-onetime" {
				var r rhandlers.JoinRoomOneTimeRes

				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				if len(c.RoomQueue) == 0 {
					continue
				}

				// get first room in queue
				first := c.RoomQueue[0]
				first.Room = &rstructs.Room{
					ID:              r.Data.RoomID,
					Name:            r.Data.RoomName,
					OneTimeJoinKeys: []string{},
					Users:           map[string]*rstructs.User{},
				}

				c.Rooms[r.Data.RoomID] = first
				c.RoomQueue = c.RoomQueue[1:]

				first.call("joined", map[string]string{
					"room_id":   first.Room.ID,
					"room_name": first.Room.Name,
				})
			} else if cmd.Cmd == "join-room-forever" {
				var r rhandlers.JoinRoomForeverRes

				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				if len(c.RoomQueue) == 0 {
					continue
				}

				// get first room in queue
				first := c.RoomQueue[0]
				first.Room = &rstructs.Room{
					ID:              r.Data.RoomID,
					Name:            r.Data.RoomName,
					OneTimeJoinKeys: []string{},
					Users:           map[string]*rstructs.User{},
				}

				c.Rooms[r.Data.RoomID] = first
				c.RoomQueue = c.RoomQueue[1:]

				first.call("joined", map[string]string{
					"room_id":   first.Room.ID,
					"room_name": first.Room.Name,
				})
			}
		}
	}()

	for {
		select {
		case <-stopRunning:
			return
		}
	}
}

func (c *Client) JoinRoom(keyType, username, roomID, key string) (*CRoom, error) {
	if keyType == "onetime" {
		req := rhandlers.JoinRoomOneTimeReq{
			Cmd: "join-room-onetime",
			Data: rhandlers.JoinRoomOneTimeData{
				RoomID:     roomID,
				OneTimeKey: key,
				Username:   username,
			},
		}

		cr := &CRoom{
			Client:       c,
			IsOwner:      false,
			ReqUsername:  username,
			Room:         nil,
			MessageQueue: []string{},
			callbacks:    map[string]interface{}{},
		}

		// add room to queue...it gets updated when the server lets us know we created it successfully
		c.RoomQueue = append(c.RoomQueue, cr)

		err := c.Conn.WriteJSON(req)
		if err != nil {
			log.Fatalln("joinRoomOnetime", err)
		}

		return cr, nil
	} else if keyType == "forever" {
		req := rhandlers.JoinRoomForeverReq{
			Cmd: "join-room-forever",
			Data: rhandlers.JoinRoomForeverReqData{
				RoomID:     roomID,
				ForeverKey: key,
				Username:   username,
			},
		}

		cr := &CRoom{
			Client:       c,
			IsOwner:      false,
			ReqUsername:  username,
			Room:         nil,
			MessageQueue: []string{},
			callbacks:    map[string]interface{}{},
		}

		// add room to queue...it gets updated when the server lets us know we created it successfully
		c.RoomQueue = append(c.RoomQueue, cr)

		err := c.Conn.WriteJSON(req)
		if err != nil {
			log.Fatalln("joinRoomForever", err)
		}

		return cr, nil
	}

	return nil, errors.New("invalid keyType")
}

func (c *Client) CreateRoom(name, username string) *CRoom {
	req := rhandlers.CreateRoomRequest{
		Cmd: "create-room",
		Data: rhandlers.CreateRoomRequestData{
			Name:     name,
			Username: username,
		},
	}

	cr := &CRoom{
		Client:       c,
		IsOwner:      true,
		ReqName:      name,
		ReqUsername:  username,
		Room:         nil,
		MessageQueue: []string{},
		callbacks:    map[string]interface{}{},
	}

	// add room to queue...it gets updated when the server lets us know we created it successfully
	c.RoomQueue = append(c.RoomQueue, cr)

	err := c.Conn.WriteJSON(req)
	if err != nil {
		log.Fatalln("createRoom", err)
	}

	return cr
}
