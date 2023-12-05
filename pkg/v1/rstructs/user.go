package rstructs

import "github.com/gorilla/websocket"

type User struct {
	Conn     *websocket.Conn
	UserID   string
	Username string
}
