package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/maximekuhn/partage/internal/core/entity"
)

// AuthenticatedMw is a middleware that redirects the user
// to /login if he isn't authenticatd when performing a request.
// This middleware can be used to restrict some routes only to authenticated
// users. It does not check any permission(s)/role(s), only the authentication.
// Important: this middleware must be placed between the [AuthMw] and the next handler
// as it uses the [AuthMw] to get the user information.
type AuthenticatedMw struct{}

func NewAuthenticatedMw() *AuthenticatedMw {
	return &AuthenticatedMw{}
}

func (a *AuthenticatedMw) AuthenticatedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from authenticated middleware")

		authmwdata, ok := r.Context().Value(AuthDatacontextKey).(AuthMwData)

		var user *entity.User
		if ok && authmwdata.Authenticated {
			user = authmwdata.User
		}

		if user == nil {
			errMsg := url.QueryEscape("You must first login to use this feature.")
			redirectUrl := fmt.Sprintf("/login?err_msg=%s", errMsg)
			http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
