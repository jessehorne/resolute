package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/jessehorne/resolute/pkg/v1/client"
	"os"
)

func main() {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	c, err := client.NewClient("/v1", "127.0.0.1:5656", tlsConf)
	if err != nil {
		panic(err)
	}

	roomID := "VGvpY6GLBds19owIMU3kkZUThnz9Rgv1"
	roomKey := "GTd5ugAgwvQ6x0LLWMxMg5WhCS6sa2eD"
	testRoom, err := c.JoinRoom("onetime", "bob", roomID, roomKey)
	if err != nil {
		panic(err)
	}

	testRoom.On("joined", func(roomID, roomName string) {
		fmt.Println("MY ID:", c.UserID)
		fmt.Println("[SUCCESS] joined room", roomID, roomName)

		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				text, _ := reader.ReadString('\n')
				if err := testRoom.SendMessage(text); err != nil {
					fmt.Println(err)
				}
				fmt.Println("SENDING: ", text)
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
