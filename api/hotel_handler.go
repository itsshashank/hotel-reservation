package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types"

	"github.com/itsshashank/hotel-reservation/types/errors"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params types.HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	var filter db.Map
	if params.Rating != 0 {
		filter = db.Map{
			"rating": params.Rating,
		}
	}
	hotels, err := h.store.Hotel.List(c.Context(), filter, &params.Pagination)
	if err != nil {
		return errors.ErrNotResourceNotFound("hotels")
	}
	resp := types.ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(params.Page),
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := db.Map{"hotelID": id}
	rooms, err := h.store.Room.List(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
