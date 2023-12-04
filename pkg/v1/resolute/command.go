package resolute

import (
	"encoding/json"
)

const (
	CommandTypeCreateRoom        = "create-room"
	CommandTypeGetRoomOneTimeKey = "room-key-onetime"
	CommandTypeGetRoomForeverKey = "room-key-forever"
	CommandTypeJoinRoomOneTime   = "join-room-onetime"
	CommandTypeJoinRoomForever   = "join-room-forever"
	CommandTypeResetRoomKeys     = "reset-room-keys"

	CommandTypeSendMessage = "send-message"
)

type Command struct {
	Cmd string `json="cmd"`
}

func NewCommandFromJSON(d []byte) (*Command, error) {
	var m Command
	err := json.Unmarshal(d, &m)
	if err != nil {
		return &m, err
	}

	return &m, nil
}
