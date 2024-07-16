package middleware

import (
	"context"
	goErrors "errors"
	"net/http"
	"os"
	"strings"

	"github.com/cinema-booker/internal/constants"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/jwt"
)

func IsAuth(handlerFunc errors.ErrorHandler, store user.UserStore) errors.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		token := r.Header.Get("Authorization")
		if token == "" {
			return errors.CustomError{
				Key: errors.Unauthorized,
				Err: goErrors.New("token is required"),
			}
		}

		tokenParts := strings.Split(token, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return errors.CustomError{
				Key: errors.Unauthorized,
				Err: goErrors.New("invalid token"),
			}
		}

		userId, err := jwt.GetTokenUserId(tokenParts[1], os.Getenv("JWT_SECRET"))
		if err != nil {
			return errors.CustomError{
				Key: errors.Unauthorized,
				Err: goErrors.New("invalid token"),
			}
		}

		user, err := store.FindById(userId)
		if err != nil {
			return err
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.UserIDKey, user.Id)
		ctx = context.WithValue(ctx, constants.UserRoleKey, user.Role)
		r = r.WithContext(ctx)

		return handlerFunc(w, r)
	}
}
