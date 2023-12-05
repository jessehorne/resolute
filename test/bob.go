package main

import (
	"fmt"
	"time"

	"github.com/jessehorne/resolute/pkg/v1/client"
)

func main() {
	c := client.NewClient("/v1", "127.0.0.1:5656")

	roomID := "iL6oB1RBw1AK1Odg6UoU397eDozV3fek"
	roomKey := "kQQ8ZKpz0PmoqQQbGILvgoMc72etKIbI"
	testRoom, err := c.JoinRoom("forever", "bob", roomID, roomKey)
	if err != nil {
		panic(err)
	}

	testRoom.On("joined", func(roomID, roomName string) {
		fmt.Println("[SUCCESS] joined room", roomID, roomName)

		go func() {
			for {
				testRoom.SendMessage("hello world")
				time.Sleep(5 * time.Second)
			}
		}()
		testRoom.On("send-message", func(roomID, userID, username, content string) {
			fmt.Println(fmt.Sprintf("[MESSAGE] RoomID: %s | UserID: %s | Username: %s | Content: %s", roomID, userID, username, content))
		})
	})

	fmt.Println("connecting to server on port 5656")
	c.Listen()
}
