package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/hotel-reservation/api/middleware"
	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types/errors"
)

type Server struct {
	app           *fiber.App
	listenAddress string
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var (
		apiError errors.ApiError
		ok       bool
	)

	if apiError, ok = err.(errors.ApiError); !ok {
		apiError = errors.NewError(http.StatusInternalServerError, err.Error())

	}

	return apiError.JSONResponse(c)
}

func NewServer(s *db.Store, ladd string) *Server {
	var (
		app = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})

		authHandler    = NewAuthHandler(s.User)
		userHandler    = NewUserHandler(s.User)
		hotelHandler   = NewHotelHandler(s)
		roomHandler    = NewRoomHandler(s)
		bookingHandler = NewBookingHandler(s)
		open           = app.Group("/api/auth")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(s.User))
		admin          = app.Group("/admin", middleware.AdminAuth)
		userRouter     = apiv1.Group("/users")
		hotelRouter    = apiv1.Group("/hotels")
		roomRouter     = apiv1.Group("/rooms")
		bookingRouter  = apiv1.Group("/bookings")
	)

	open.Post("/", authHandler.HandleAuthenticate)
	open.Post("/api/v1/users", userHandler.HandlePostUser)

	admin.Get("/users", userHandler.HandleGetUsers)
	admin.Get("/users/:id", userHandler.HandleGetUser)
	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	userRouter.Get("/", userHandler.HandleGetUser)
	userRouter.Patch("/", userHandler.HandlePatchUser)
	userRouter.Delete("/:id", userHandler.HandleDeleteUser)

	hotelRouter.Get("/", hotelHandler.HandleGetHotels)
	hotelRouter.Get("/:id", hotelHandler.HandleGetHotel)
	hotelRouter.Get("/:id/rooms", hotelHandler.HandleGetRooms)

	roomRouter.Get("/", roomHandler.HandleGetRooms)
	roomRouter.Post("/:id/book", roomHandler.HandleBookRoom)

	bookingRouter.Get("/:id", bookingHandler.HandleGetBooking)
	bookingRouter.Post("/:id/cancel", bookingHandler.HandleCancelBooking)

	return &Server{
		app:           app,
		listenAddress: ladd,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(s.listenAddress)
}
