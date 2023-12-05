package rstructs

import (
	"github.com/dchest/uniuri"
)

type State struct {
	Rooms map[string]*Room
}

func NewState() *State {
	return &State{
		Rooms: map[string]*Room{},
	}
}

func (s *State) AddRoom(room *Room) {
	s.Rooms[room.ID] = room
}

func (s *State) HasRoom(id string) bool {
	_, ok := s.Rooms[id]
	return ok
}

func (s *State) CreateOneTimeRoomKey(id string) string {
	r, ok := s.Rooms[id]

	if !ok {
		return ""
	}

	newKey := uniuri.NewLen(32)
	r.OneTimeJoinKeys = append(r.OneTimeJoinKeys, newKey)

	return newKey
}

func (s *State) CleanupRooms(userID string) {
	for k, r := range s.Rooms {
		if r.OwnerID == userID {
			delete(s.Rooms, k)
		}
	}
}
