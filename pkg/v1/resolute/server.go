package resolute

import (
	"log"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/gorilla/websocket"
	handlers2 "github.com/jessehorne/resolute/pkg/v1/rhandlers"
	"github.com/jessehorne/resolute/pkg/v1/rstructs"
)

var State = rstructs.NewState()

func serverHandler(w http.ResponseWriter, r *http.Request) {
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
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				break
			}
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
		return handlers2.CreateRoomHandler(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeGetRoomOneTimeKey {
		return handlers2.GetRoomOneTimeKey(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeGetRoomForeverKey {
		return handlers2.GetRoomForeverKey(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeJoinRoomOneTime {
		return handlers2.JoinRoomOneTime(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeJoinRoomForever {
		return handlers2.JoinRoomForever(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeResetRoomKeys {
		return handlers2.ResetRoomKeys(State, userID, c, m)
	} else if cmd.Cmd == CommandTypeSendMessage {
		return handlers2.SendMessage(State, userID, c, m)
	}

	return nil
}

type Server struct {
	Path    string
	Host    string
	Handler func(http.ResponseWriter, *http.Request)
}

func NewServer(path, host string) *Server {
	return &Server{
		Path:    path,
		Host:    host,
		Handler: serverHandler,
	}
}

func (s *Server) Listen() error {
	http.HandleFunc(s.Path, s.Handler)
	return http.ListenAndServe(s.Host, nil)
}
