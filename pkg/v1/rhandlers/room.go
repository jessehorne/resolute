package rhandlers

type RoomIDReq struct {
	RoomID string `json:"room_id"`
}

type ResponseError struct {
	Cmd  string `json:"cmd"`
	Data map[string]string
}
