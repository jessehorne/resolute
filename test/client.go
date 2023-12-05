package main

import (
	"fmt"

	"github.com/jessehorne/resolute/pkg/v1/client"
)

func main() {
	c := client.NewClient("/v1", "127.0.0.1:5656")

	testRoom := c.CreateRoom("test room", "bob")

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
	})

	fmt.Println("connecting to server on port 5656")
	c.Listen()
}
