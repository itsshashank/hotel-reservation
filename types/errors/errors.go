package errors

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ApiError interface {
	error
	JSONResponse(ctx *fiber.Ctx) error
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func (e Error) JSONResponse(ctx *fiber.Ctx) error {
	return ctx.Status(e.Code).JSON(e.Err)
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}

func ErrNotResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + " resource not found",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func InvalidCredentials() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "invalid credentials",
	}
}

func ErrRoomBooked(rid string) Error {
	return Error{
		Code: http.StatusNotAcceptable,
		Err:  fmt.Sprintf("room %s already booked", rid),
	}
}
