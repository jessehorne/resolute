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

	roomID := "08qXiXm4J5Ad600j26uvxpbRmLvkDaJ1"
	roomKey := "xHqSrGvcOV51LQt3i2zxKcH28Ge3meFd"
	testRoom, err := c.JoinRoom("onetime", "bob", roomID, roomKey)
	if err != nil {
		panic(err)
	}

	testRoom.On("joined", func(roomID, roomName string) {
		fmt.Println("MY ID:", c.UserID)
		fmt.Println("[SUCCESS] joined room", roomID, roomName)

		go func() {
			for {
				testRoom.SendMessage("hello world from bob")
				time.Sleep(5 * time.Second)
			}
		}()
		testRoom.On("send-message", func(roomID, userID, username, content string) {
			fmt.Println(fmt.Sprintf("[MESSAGE] RoomID: %s | Username: %s | Content: %s",
				roomID, username, content))
		})
	})

	fmt.Println("connecting to server on port 5656")
	c.Listen()
}
