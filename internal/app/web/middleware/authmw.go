package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/maximekuhn/partage/internal/auth"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/query"
)

type AuthContextKey string

const (
	AuthDatacontextKey AuthContextKey = "partage-auth-mw-authData"
	cookieName         string         = "authToken"
)

// AuthMw is an authentication middleware that
// inject [AuthMwData] into the request context.
// This middleware doesn't return any error if the user isn't authenticated.
type AuthMw struct {
	authSvc            *auth.AuthService
	getUserByIDHandler *query.GetUserByIDQueryHandler
}

func NewAuthMw(authSvc *auth.AuthService, getUserByIDHandler *query.GetUserByIDQueryHandler) *AuthMw {
	return &AuthMw{authSvc, getUserByIDHandler}
}

func (a *AuthMw) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from auth middleware")

		var authData AuthMwData

		cookie := r.Header.Get("Cookie")
		parts := strings.Split(cookie, fmt.Sprintf("%s=", cookieName))
		if len(parts) != 2 {
			fmt.Println("could not find auth token in cookie")
			authData.Authenticated = false
			ctx := context.WithValue(r.Context(), AuthDatacontextKey, authData)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := parts[1]

		userID, err := a.authSvc.VerifyToken(tokenStr)
		if err != nil {
			fmt.Println("could not verify token")
			authData.Authenticated = false
			ctx := context.WithValue(r.Context(), AuthDatacontextKey, authData)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		user, found, err := a.getUserByIDHandler.Handle(r.Context(), query.GetUserByIDQuery{
			ID: *userID,
		})
		if err != nil || !found {
			fmt.Println("could not get user or user not found")
			authData.Authenticated = false
			ctx := context.WithValue(r.Context(), AuthDatacontextKey, authData)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		authData.User = user
		authData.Authenticated = true

		ctx := context.WithValue(r.Context(), AuthDatacontextKey, authData)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

type AuthMwData struct {
	// The current user, guaranteed to be non-nil
	// if Authenticated is true.
	User          *entity.User
	Authenticated bool
}
