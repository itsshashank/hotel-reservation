package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types"
	"github.com/itsshashank/hotel-reservation/types/errors"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.List(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {

	var (
		params types.BookRoomParams
		roomID = c.Params("id")
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	fromDate, err := time.Parse(time.DateOnly, params.FromDate)
	if err != nil {
		return err
	}
	tillDate, err := time.Parse(time.DateOnly, params.TillDate)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return errors.NewError(http.StatusInternalServerError, "internal server error")
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   fromDate,
		TillDate:   tillDate,
		NumPersons: params.NumPersons,
	}

	if err := booking.Validate(); err != nil {
		return err
	}

	ok, err = h.store.Room.IsAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return errors.ErrRoomBooked(c.Params("id"))
	}

	inserted, err := h.store.Booking.Insert(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}
