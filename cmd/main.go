package main

import (
	"log"
	"os"

	"github.com/itsshashank/hotel-reservation/api"
	"github.com/itsshashank/hotel-reservation/db"
)

var dburi, dbname, listenAddr string

func main() {
	store := db.Store{
		User:    db.NewMogoUserStore(dburi, dbname, db.UserColl),
		Hotel:   db.NewMongoHotelStore(dburi, dbname, db.HotelColl),
		Room:    db.NewMongoRoomStore(dburi, dbname, db.RoomColl),
		Booking: db.NewMongoBookingStore(dburi, dbname, db.BookingColl),
	}

	server := api.NewServer(&store, listenAddr)
	server.Start()
}

func init() {
	dburi = os.Getenv("MONGODB_URI")
	if dburi == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	dbname = os.Getenv("DB_NAME")

	listenAddr = os.Getenv("HTTP_LISTEN_ADDRESS")
	if listenAddr == "" {
		listenAddr = ":5000"
	}
}
