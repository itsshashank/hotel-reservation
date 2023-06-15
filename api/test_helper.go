package api

import (
	"context"
	"os"
	"testing"

	"github.com/itsshashank/hotel-reservation/db"
)

var testdbName = "test_db"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	dburi := os.Getenv("MONGODB_URI")

	dbColl := "test_coll"
	userStore := db.NewMogoUserStore(dburi, testdbName, dbColl)
	return &testdb{
		UserStore: userStore,
	}
}
