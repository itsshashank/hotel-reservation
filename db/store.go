package db

import "context"

type Map map[string]any

type Storer interface {
	Drop(context.Context) error
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
