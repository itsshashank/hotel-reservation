package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types"
	"github.com/itsshashank/hotel-reservation/types/errors"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(s db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: s,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {

	var params types.CreateUserParams

	if err := c.BodyParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	if error := params.Validate(); len(error) > 0 {
		return c.JSON(error)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.Insert(c.Context(), user)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(insertedUser)
}

func (h *UserHandler) HandlePatchUser(c *fiber.Ctx) error {
	var params types.UpdateUserParams

	user, err := types.GetAuthUser(c)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	if err := h.userStore.Update(c.Context(), user.ID, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": user.ID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.Delete(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userID})
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {
		user, err := types.GetAuthUser(c)
		if err != nil {
			return err
		}
		return c.JSON(user)
	}

	user, err := h.userStore.GetByID(c.Context(), id)
	if err != nil {
		return errors.NewError(http.StatusNotFound, err.Error())
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.List(c.Context())
	if err != nil {
		return errors.ErrNotResourceNotFound("users")
	}
	return c.JSON(users)
}
