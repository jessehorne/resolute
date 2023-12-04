package structs

import "github.com/dchest/uniuri"

type Room struct {
	ID              string   `json:"room_id"`
	OwnerID         string   `json:"owner_id"`
	Name            string   `json:"name"`
	AdminKey        string   `json:"admin_key"`
	OneTimeJoinKeys []string `json:"one_time_join_keys"`
	ForeverJoinKey  string   `json:"forever_join_key"`
	Users           []string `json:"users"`
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
		Users:           []string{},
	}
}

func (r *Room) HasUser(userID string) bool {
	for _, k := range r.Users {
		if k == userID {
			return true
		}
	}

	return false
}

// AddUser adds a userID to a room if it isn't already there. This allows the user to receive messages
// from the room.
func (r *Room) AddUser(userID string) {
	if !r.HasUser(userID) {
		r.Users = append(r.Users, userID)
	}
}
