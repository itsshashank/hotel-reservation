package api

import (
	e "errors"

	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types"
	"github.com/itsshashank/hotel-reservation/types/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User    *types.User `json:"user"`
	Token   string      `json:"token"`
	Expires string      `json:"expires"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.userStore.GetByEmail(c.Context(), params.Email)
	if err != nil {
		if e.Is(err, mongo.ErrNoDocuments) {
			return errors.InvalidCredentials()
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return errors.InvalidCredentials()
	}

	token, expires := user.CreateToken()

	resp := AuthResponse{
		User:    user,
		Token:   token,
		Expires: expires,
	}
	return c.JSON(resp)
}
