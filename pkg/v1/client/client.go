package client

import (
	"crypto/rsa"
	"crypto/tls"
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
	"github.com/jessehorne/resolute/pkg/v1/util"
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
	Rooms     []*CRoom // rooms the user is in. key is the room id
	RoomQueue []*CRoom

	Conn *websocket.Conn

	PrivateKey *rsa.PrivateKey
}

func NewClient(path, host string, tlsConf *tls.Config) (*Client, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: path}

	var c *websocket.Conn

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = tlsConf

	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// generate keys that will last the lifetime of this users connection
	key, err := util.GenerateAsymmetricKey()
	if err != nil {
		return nil, err
	}

	newClient := &Client{
		Host:       host,
		Path:       path,
		Rooms:      []*CRoom{},
		RoomQueue:  []*CRoom{},
		Conn:       c,
		PrivateKey: key,
	}

	return newClient, nil
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

				// get room and update values
				for _, ro := range c.Rooms {
					if ro.ClientRoomID == r.ClientRoomID {
						ro.Room = &rstructs.Room{
							ID:              r.Room.ID,
							OwnerID:         r.Room.OwnerID,
							Name:            r.Room.Name,
							OneTimeJoinKeys: []string{},
							ForeverJoinKey:  "",
							Users:           map[string]*rstructs.User{},
						}

						// update client id
						c.UserID = r.Room.OwnerID
						ro.call("created", r.Room.ID)

						break
					}
				}
			} else if cmd.Cmd == "room-key-onetime" {
				var r rhandlers.GetRoomOneTimeKeyResponse
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				for _, ro := range c.Rooms {
					if ro.Room.ID == r.Data.RoomID {
						ro.call("key-onetime", map[string]string{
							"room_id": r.Data.RoomID,
							"key":     r.Data.OneTimeKey,
						})

						break
					}
				}
			} else if cmd.Cmd == "room-key-forever" {
				var r rhandlers.GetRoomForeverKeyResponse
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				for _, ro := range c.Rooms {
					if ro.Room.ID == r.Data.RoomID {
						ro.call("key-forever", map[string]string{
							"room_id": r.Data.RoomID,
							"key":     r.Data.ForeverJoinKey,
						})

						break
					}
				}
			} else if cmd.Cmd == "send-message" {
				var r rhandlers.SendMessageRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
					return
				}

				d, err := util.DecryptMessage(c.PrivateKey, r.Data.Content)
				if err != nil {
					return
				}

				for _, ro := range c.Rooms {
					if ro.Room.ID == r.Data.RoomID {
						ro.call("send-message", map[string]string{
							"room_id":  r.Data.RoomID,
							"user_id":  r.Data.UserID,
							"username": r.Data.Username,
							"content":  d,
						})

						break
					}
				}
			} else if cmd.Cmd == "user-joined-onetime" {
				var r rhandlers.JoinRoomOneTimeRes

				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				for _, ro := range c.Rooms {
					if ro.Room.ID == r.Data.RoomID {

						// add this user to list of users in room
						ro.AddJoinedUser(r.Data.UserID, r.Data.Username, r.Data.PublicKey)

						ro.call("user-joined", map[string]string{
							"room_id":   r.Data.RoomID,
							"room_name": r.Data.RoomName,
							"user_id":   r.Data.UserID,
							"username":  r.Data.Username,
							"key_type":  "forever",
						})

						break
					}
				}
			} else if cmd.Cmd == "user-joined-forever" {
				var r rhandlers.JoinRoomForeverRes

				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				for _, ro := range c.Rooms {
					if ro.Room.ID == r.Data.RoomID {

						// add this user to list of users in room
						ro.AddJoinedUser(r.Data.UserID, r.Data.Username, r.Data.PublicKey)

						ro.call("user-joined", map[string]string{
							"room_id":   r.Data.RoomID,
							"room_name": r.Data.RoomName,
							"user_id":   r.Data.UserID,
							"username":  r.Data.Username,
							"key_type":  "forever",
						})

						break
					}
				}
			} else if cmd.Cmd == "joined" {
				var r rhandlers.YouJoinedRoomRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					// TODO
				}

				for _, ro := range c.Rooms {
					if ro.Room.ID == r.Data.RoomID {
						c.UserID = r.Data.UserID
						ro.UpdateUsers(r.Data.Users)
						ro.call("joined", map[string]string{
							"room_id":   r.Data.RoomID,
							"room_name": r.Data.RoomName,
						})

						break
					}
				}
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
	pubKeyString := util.PublicKeyToString(&c.PrivateKey.PublicKey)

	if keyType == "onetime" {
		cr := &CRoom{
			Client:      c,
			IsOwner:     false,
			ReqUsername: username,
			Room: &rstructs.Room{
				ID: roomID,
			},
			MessageQueue: []string{},
			callbacks:    map[string]interface{}{},
		}
		c.Rooms = append(c.Rooms, cr)

		req := rhandlers.JoinRoomOneTimeReq{
			Cmd: "join-room-onetime",
			Data: rhandlers.JoinRoomOneTimeReqData{
				RoomID:     roomID,
				OneTimeKey: key,
				Username:   username,
				PublicKey:  pubKeyString,
			},
		}
		err := c.Conn.WriteJSON(req)
		if err != nil {
			log.Fatalln("joinRoomOnetime", err)
		}

		return cr, nil
	} else if keyType == "forever" {
		cr := &CRoom{
			Client:      c,
			IsOwner:     false,
			ReqUsername: username,
			Room: &rstructs.Room{
				ID: roomID,
			},
			MessageQueue: []string{},
			callbacks:    map[string]interface{}{},
		}
		c.Rooms = append(c.Rooms, cr)

		req := rhandlers.JoinRoomForeverReq{
			Cmd: "join-room-forever",
			Data: rhandlers.JoinRoomForeverReqData{
				RoomID:     roomID,
				ForeverKey: key,
				Username:   username,
				PublicKey:  pubKeyString,
			},
		}
		err := c.Conn.WriteJSON(req)
		if err != nil {
			log.Fatalln("joinRoomForever", err)
		}

		return cr, nil
	}

	return nil, errors.New("invalid keyType")
}

func (c *Client) CreateRoom(name, username, clientRoomID string) *CRoom {
	pubKeyString := util.PublicKeyToString(&c.PrivateKey.PublicKey)

	req := rhandlers.CreateRoomRequest{
		Cmd: "create-room",
		Data: rhandlers.CreateRoomRequestData{
			Name:            name,
			Username:        username,
			ClientRoomID:    clientRoomID,
			PublicKeyString: pubKeyString,
		},
	}

	cr := &CRoom{
		ClientRoomID: clientRoomID,
		Client:       c,
		IsOwner:      true,
		ReqName:      name,
		ReqUsername:  username,
		Room:         nil,
		MessageQueue: []string{},
		callbacks:    map[string]interface{}{},
	}

	c.Rooms = append(c.Rooms, cr)

	err := c.Conn.WriteJSON(req)
	if err != nil {
		log.Fatalln("createRoom", err)
	}

	return cr
}
