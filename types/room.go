package types

type RoomType string

const (
	Single  RoomType = "single"
	Double  RoomType = "double"
	SeaSide RoomType = "seaside"
	Deluxe  RoomType = "deluxe"
)

type Room struct {
	ID      string   `bson:"_id,omitempty" json:"id,omitempty"`
	Type    RoomType `bson:"type" json:"type"`
	Price   float64  `bson:"price" json:"price"`
	HotelID string   `bson:"hotelID" json:"hotelID"`
}

type BookRoomParams struct {
	FromDate   string `json:"fromDate"`
	TillDate   string `json:"tillDate"`
	NumPersons int    `json:"numPersons"`
}
