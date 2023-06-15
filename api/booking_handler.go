package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types"
	"github.com/itsshashank/hotel-reservation/types/errors"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.List(c.Context(), db.Map{})
	if err != nil {
		return errors.ErrNotResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return errors.ErrNotResourceNotFound("booking")
	}
	user, err := types.GetAuthUser(c)
	if err != nil {
		return errors.ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return errors.ErrUnAuthorized()
	}
	if err := h.store.Booking.Update(c.Context(), c.Params("id"), db.Map{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(map[string]string{"Cancled Booking": id})
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return errors.ErrNotResourceNotFound("booking")
	}
	user, err := types.GetAuthUser(c)
	if err != nil {
		return errors.ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return errors.ErrUnAuthorized()
	}
	return c.JSON(booking)
}
