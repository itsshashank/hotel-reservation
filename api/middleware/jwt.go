package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"github.com/itsshashank/hotel-reservation/db"
	"github.com/itsshashank/hotel-reservation/types/errors"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("token not present in the header")
			return errors.ErrUnAuthorized()
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return errors.NewError(http.StatusUnauthorized, "token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetByID(ctx.Context(), userID)
		if err != nil {
			fmt.Println(err)
			return errors.ErrUnAuthorized()
		}
		// Set the current authenticated user to the context.
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, errors.ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, errors.ErrUnAuthorized()
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, errors.ErrUnAuthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrUnAuthorized()
	}
	return claims, nil
}
