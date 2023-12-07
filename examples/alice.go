package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/jessehorne/resolute/pkg/v1/client"
)

func main() {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	c := client.NewClient("/v1", "127.0.0.1:5656", tlsConf)

	testRoom := c.CreateRoom("test room", "alice")

	testRoom.On("created", func() {
		fmt.Println("[SUCCESS] room created")

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
