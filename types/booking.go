package types

import (
	"fmt"
	"time"
)

type Booking struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     string    `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID     string    `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons int       `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	FromDate   time.Time `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate   time.Time `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	Canceled   bool      `bson:"canceled" json:"canceled"`
}

func (b Booking) Validate() error {
	now := time.Now()
	if now.After(b.FromDate) || now.After(b.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	return nil
}
