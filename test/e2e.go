package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/handlers"
	"github.com/jessehorne/resolute/resolute"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:5656", Path: "/v1"}
	fmt.Println("Connecting to 5656")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			cmd, err := resolute.NewCommandFromJSON(message)
			if err != nil {
				log.Fatalln(err)
			}

			if cmd.Cmd == "create-room" {
				var r handlers.CreateRoomResponse
				err := json.Unmarshal(message, &r)
				if err != nil {
					log.Fatalln("create-room:", err)
				}
				fmt.Println("[SUCCESS] created room")

				onetimeKey := fmt.Sprintf(`{"cmd": "room-key-onetime", "data": { "room_id": "%s" }}`,
					r.Room.ID)
				err = c.WriteMessage(websocket.TextMessage, []byte(onetimeKey))
				if err != nil {
					log.Fatalln("room-key-onetime", err)
				}

				foreverKey := fmt.Sprintf(`{"cmd": "room-key-forever", "data": { "room_id": "%s" }}`,
					r.Room.ID)
				err = c.WriteMessage(websocket.TextMessage, []byte(foreverKey))
				if err != nil {
					log.Fatalln("room-key-forever", err)
				}

				// send test message
				sendMsg := fmt.Sprintf(`{"cmd": "send-message", "data": { "room_id": "%s", "content": "test message" }}`,
					r.Room.ID)
				err = c.WriteMessage(websocket.TextMessage, []byte(sendMsg))
				if err != nil {
					log.Fatalln("sendMessage", err)
				}
			} else if cmd.Cmd == "room-key-onetime" {
				var r handlers.GetRoomOneTimeKeyResponse
				err := json.Unmarshal(message, &r)
				if err != nil {
					log.Fatalln("room-key-onetime:", err)
				}

				fmt.Println("[SUCCESS] get onetime join key", r.Data.RoomID, r.Data.OneTimeKey)

				// attempt to join room with onetime key
				joinRoomOneTime := fmt.Sprintf(`{"cmd": "join-room-onetime", "data": { "room_id": "%s", "one_time_key": "%s" }}`,
					r.Data.RoomID, r.Data.OneTimeKey)
				err = c.WriteMessage(websocket.TextMessage, []byte(joinRoomOneTime))
				if err != nil {
					log.Fatalln("joinRoomOneTime", err)
				}
			} else if cmd.Cmd == "room-key-forever" {
				var r handlers.RoomForeverKeyGetRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					log.Fatalln("room-key-forever:", err)
				}

				fmt.Println("[SUCCESS] get forever join key", r.Data.RoomID, r.Data.ForeverJoinKey)

				// attempt to join room with forever key
				joinRoomForever := fmt.Sprintf(`{"cmd": "join-room-forever", "data": { "room_id": "%s", "forever_key": "%s" }}`,
					r.Data.RoomID, r.Data.ForeverJoinKey)
				err = c.WriteMessage(websocket.TextMessage, []byte(joinRoomForever))
				if err != nil {
					log.Fatalln("joinRoomForever", err)
				}
			} else if cmd.Cmd == "join-room-onetime" {
				var r handlers.JoinRoomOneTimeRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					log.Fatalln("join-room-onetime:", err)
				}

				fmt.Println("[SUCCESS] join room using onetime key", r.Data.RoomID)
			} else if cmd.Cmd == "join-room-forever" {
				var r handlers.JoinRoomForeverRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					log.Fatalln("join-room-forever:", err)
				}

				fmt.Println("[SUCCESS] join room using forever key", r.Data.RoomID)

				// reset all keys
				resetKeys := fmt.Sprintf(`{"cmd": "reset-room-keys", "data": { "room_id": "%s" }}`,
					r.Data.RoomID)
				err = c.WriteMessage(websocket.TextMessage, []byte(resetKeys))
				if err != nil {
					log.Fatalln("joinRoomForever", err)
				}
			} else if cmd.Cmd == "reset-room-keys" {
				var r handlers.ResetRoomKeysRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					log.Fatalln("reset-room-keys:", err)
				}

				fmt.Println("[SUCCESS] reset room keys", r.Data.RoomID)
			} else if cmd.Cmd == "send-message" {
				var r handlers.SendMessageRes
				err := json.Unmarshal(message, &r)
				if err != nil {
					log.Fatalln("send-message:", err)
				}

				fmt.Println("[SUCCESS] sent message", r.Data.RoomID, r.Data.Content, r.Data.From)
			}

			var r handlers.GetRoomOneTimeKeyResponse
			err = json.Unmarshal(message, &r)
			if err != nil {
				log.Println("room key test:", err)
			}

			//log.Println("recv:", string(message))
		}
	}()

	// TESTS

	// Create room
	createRoom := `{"cmd": "create-room", "data": { "name": "test room 123" }}`
	err = c.WriteMessage(websocket.TextMessage, []byte(createRoom))
	if err != nil {
		log.Fatalln("createRoom", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			return
		}
	}
}
