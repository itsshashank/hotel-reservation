package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/hotel-reservation/types"
	"github.com/itsshashank/hotel-reservation/types/errors"
)

func AdminAuth(ctx *fiber.Ctx) error {
	if user, ok := ctx.Context().UserValue("user").(*types.User); ok {
		if user.IsAdmin {
			return ctx.Next()
		}
	}
	return errors.ErrUnAuthorized().JSONResponse(ctx)
}
