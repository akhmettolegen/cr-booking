package entity

import "time"

type Reservation struct {
	Id        string    `json:"id"`
	RoomId    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
