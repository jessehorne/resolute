package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/jessehorne/resolute/pkg/v1/client"
)

func main() {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	c, err := client.NewClient("/v1", "127.0.0.1:5656", tlsConf)
	if err != nil {
		panic(err)
	}

	testRoom := c.CreateRoom("test room", "alice", "whatever room id")

	testRoom.On("created", func(roomID string) {
		fmt.Println("MY ID:", c.UserID)
		fmt.Println("[SUCCESS] room created | RoomID: ", roomID)

		testRoom.GetKey("onetime")
		testRoom.On("key-onetime", func(roomID, key string) {
			fmt.Println("[SUCCESS] got onetime key:", key)
		})

		testRoom.GetKey("forever")
		testRoom.On("key-forever", func(roomID, key string) {
			fmt.Println("[SUCCESS] got forever key:", key)
		})

		go func() {
			for {
				testRoom.SendMessage("hello world from alice")
				time.Sleep(5 * time.Second)
			}
		}()
		testRoom.On("send-message", func(roomID, userID, username, content string) {
			fmt.Println(fmt.Sprintf("[MESSAGE] RoomID: %s | Username: %s | Content: %s",
				roomID, username, content))
		})

		testRoom.On("user-joined", func(roomID, roomName, userID, username, keyType string) {
			fmt.Println("[USER JOINED] ", roomID, roomName, userID, username, keyType)
		})
	})

	fmt.Println("connecting to server on port 5656")
	c.Listen()
}
