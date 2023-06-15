package db

import (
	"context"

	"github.com/itsshashank/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const RoomColl = "rooms"

type RoomStore interface {
	Insert(context.Context, *types.Room) (*types.Room, error)
	List(context.Context, Map) ([]*types.Room, error)
	IsAvailableForBooking(context.Context, string, types.BookRoomParams) (bool, error)
}

type MongoRoomStore struct {
	client    *mongo.Client
	coll      *mongo.Collection
	hotelColl *mongo.Collection
}

func NewMongoRoomStore(dburi string, dbname string, collection string) *MongoRoomStore {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}
	return &MongoRoomStore{
		client:    client,
		coll:      client.Database(dbname).Collection(collection),
		hotelColl: client.Database(dbname).Collection(HotelColl),
	}
}

func (s *MongoRoomStore) Insert(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID).Hex()

	hid, err := primitive.ObjectIDFromHex(room.HotelID)
	if err != nil {
		return nil, err
	}
	// update the hotel with this room id
	filter := Map{"_id": hid}
	update := Map{"$push": bson.M{"rooms": room.ID}}
	if _, err := s.hotelColl.UpdateOne(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *MongoRoomStore) List(ctx context.Context, filter Map) ([]*types.Room, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) IsAvailableForBooking(ctx context.Context, roomID string, params types.BookRoomParams) (bool, error) {
	rid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return false, err
	}
	where := Map{
		"roomID": rid,
		"fromDate": Map{
			"$gte": params.FromDate,
		},
		"tillDate": Map{
			"$lte": params.TillDate,
		},
	}
	bookings, err := s.List(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
