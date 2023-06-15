package fixtures

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types"
)

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.Insert(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddRoom(store *db.Store, roomtype types.RoomType, price float64, hotelID string) *types.Room {
	room := &types.Room{
		Type:    roomtype,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.Insert(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddHotel(store *db.Store, name string, loc string, rating int, rooms map[types.RoomType]int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    []string{},
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range rooms {
		for j := 0; j < v; j++ {
			price := float64(rand.Intn(100)) * rand.Float64()
			AddRoom(store, k, price, insertedHotel.ID)
		}
	}
	return insertedHotel
}

func AddBooking(store *db.Store, uid, rid string, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.Insert(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
