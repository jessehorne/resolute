package rstructs

import (
	"crypto/rsa"

	"github.com/gorilla/websocket"
)

type User struct {
	Conn            *websocket.Conn
	UserID          string
	Username        string
	PublicKey       *rsa.PublicKey
	PublicKeyString string
}

type JoinedUser struct {
	UserID          string `json:"user_id"`
	Username        string `json:"username"`
	PublicKeyString string `json:"public_key_string"`
}
