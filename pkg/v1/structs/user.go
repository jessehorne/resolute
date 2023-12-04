package structs

import "github.com/gorilla/websocket"

type User struct {
	Conn   *websocket.Conn
	UserID string
}
