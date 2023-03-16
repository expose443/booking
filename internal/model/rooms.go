package model

import "time"

type Rooms struct {
	Room_id       int64
	Room_name     string
	Room_occupied int64
}

type Expires struct {
	Room_id   int64
	Check_in  time.Time
	Check_out time.Time
}

type Reservation struct {
	UserId   int64
	RoomID   int64
	CheckIn  time.Time
	CheckOut time.Time
	Author   User
	Room     Rooms
	In       string
	Out      string
}
