package rstructs

import (
	"github.com/dchest/uniuri"
)

type Room struct {
	ID              string           `json:"room_id"`
	OwnerID         string           `json:"owner_id"`
	Name            string           `json:"name"`
	AdminKey        string           `json:"admin_key"`
	OneTimeJoinKeys []string         `json:"one_time_join_keys"`
	ForeverJoinKey  string           `json:"forever_join_key"`
	Users           map[string]*User `json:"users"`
}

func NewRoom(name, ownerID string) *Room {
	id := uniuri.NewLen(32)
	adminKey := uniuri.NewLen(32)

	return &Room{
		ID:              id,
		OwnerID:         ownerID,
		Name:            name,
		AdminKey:        adminKey,
		OneTimeJoinKeys: []string{},
		ForeverJoinKey:  "",
		Users:           map[string]*User{},
	}
}

func (r *Room) HasUser(userID string) bool {
	_, ok := r.Users[userID]
	return ok
}

// AddUser adds a userID to a room if it isn't already there. This allows the user to receive messages
// from the room.
func (r *Room) AddUser(u *User) {
	r.Users[u.UserID] = u
}

type BroadcastMessageResData struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

type BroadcastMessageRes struct {
	Cmd  string                  `json:"cmd"`
	Data BroadcastMessageResData `json:"data"`
}

func (r *Room) BroadcastMessage(userID, content string) {
	u, ok := r.Users[userID]
	if !ok {
		return
	}

	for _, user := range r.Users {
		user.Conn.WriteJSON(BroadcastMessageRes{
			Cmd: "send-message",
			Data: BroadcastMessageResData{
				RoomID:   r.ID,
				UserID:   u.UserID,
				Username: u.Username,
				Content:  content,
			},
		})
	}
}
