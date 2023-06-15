package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/db/fixtures"
	"github.com/itsshashank/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var (
		dburi  = os.Getenv("MONGODB_URI")
		dbname = os.Getenv("DB_NAME")
	)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(dbname).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	store := &db.Store{
		User:    db.NewMogoUserStore(dburi, dbname, db.UserColl),
		Hotel:   db.NewMongoHotelStore(dburi, dbname, db.HotelColl),
		Room:    db.NewMongoRoomStore(dburi, dbname, db.RoomColl),
		Booking: db.NewMongoBookingStore(dburi, dbname, db.BookingColl),
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	token, _ := user.CreateToken()
	fmt.Println("james ->", token)
	admin := fixtures.AddUser(store, "admin", "admin", true)
	token, _ = admin.CreateToken()
	fmt.Println("admin ->", token)
	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, types.SeaSide, 88.44, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking.ID)

	roomtypes := []types.RoomType{types.Single, types.Double, types.Deluxe, types.SeaSide}
	for i := 0; i < 5; i++ {
		rooms := make(map[types.RoomType]int)
		name := fmt.Sprintf("hotel_%d", i)
		location := fmt.Sprintf("location_%d", i)
		randrooms := rand.Intn(10)
		for i := 0; i < randrooms; i++ {
			rooms[roomtypes[rand.Intn(len(roomtypes))]] = rand.Intn(3)
		}
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, rooms)
	}
}
