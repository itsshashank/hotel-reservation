package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/itsshashank/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const UserColl = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetByEmail(context.Context, string) (*types.User, error)
	GetByID(context.Context, string) (*types.User, error)
	List(context.Context) ([]*types.User, error)
	Insert(context.Context, *types.User) (*types.User, error)
	Delete(context.Context, string) error
	Update(ctx context.Context, ID string, params types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMogoUserStore(dburi string, dbname string, collection string) *MongoUserStore {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(collection),
	}
}

func (s *MongoUserStore) Update(ctx context.Context, ID string, params types.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	update := bson.M{"$set": params.ToBSON()}
	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	log.Println("Deleted records:", res.DeletedCount)
	return nil
}

func (s *MongoUserStore) Insert(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (s *MongoUserStore) List(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = fmt.Errorf("user: %s not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	log.Println("--- dropping user collection")
	return s.coll.Drop(ctx)
}
