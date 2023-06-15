package types

type Hotel struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Location string   `bson:"location" json:"location"`
	Rooms    []string `bson:"rooms"`
	Rating   int      `bson:"rating" json:"rating"`
}

type HotelQueryParams struct {
	Pagination
	Rating int
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}
