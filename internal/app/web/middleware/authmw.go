package middleware

import (
	"fmt"
	"net/http"

	"github.com/maximekuhn/partage/internal/auth"
	"github.com/maximekuhn/partage/internal/core/entity"
)

type AuthContextKey string

const (
	AuthDatacontextKey AuthContextKey = "partage-auth-mw-authData"
)

// AuthMw is an authentication middleware that
// inject [AuthMwData] into the request context.
// This middleware doesn't return any error if the user isn't authenticated.
type AuthMw struct {
	authSvc *auth.AuthService
}

func NewAuthMw(authSvc *auth.AuthService) *AuthMw {
	return &AuthMw{authSvc}
}

func (a *AuthMw) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello from auth middleware")
		next.ServeHTTP(w, r)
	})
}

type AuthMwData struct {
	// The current user, guaranteed to be non-nil
	// if Authenticated is true.
	User          *entity.User
	Authenticated bool
}
