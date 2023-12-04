package resolute

import (
	"log"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/gorilla/websocket"
	"github.com/jessehorne/resolute/handlers"
	"github.com/jessehorne/resolute/structs"
)

var State = structs.NewState()

func ServerHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
	}
	defer c.Close()

	userID := uniuri.NewLen(32)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		err = handleMessage(c, userID, msg)
		if err != nil {
			log.Println("handle:", err)
		}
	}
}

func handleMessage(c *websocket.Conn, userID string, m []byte) error {
	cmd, err := NewCommandFromJSON(m)
	if err != nil {
		return err
	}

	if cmd.Cmd == CommandTypeCreateRoom {
		return handlers.CreateRoomHandler(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeGetRoomOneTimeKey {
		return handlers.GetRoomOneTimeKey(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeGetRoomForeverKey {
		return handlers.GetRoomForeverKey(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeJoinRoomOneTime {
		return handlers.JoinRoomOneTime(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeJoinRoomForever {
		return handlers.JoinRoomForever(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeResetRoomKeys {
		return handlers.ResetRoomKeys(State, userID, c, m)
	}

	return nil
}
