package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/query"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

// UserInGroupMw is a middleware that checks if the user is in the
// request group.
// The group ID is get from path parameters ("id") and it returns
// a 4xx if the user isn't in the requested group.
// This middleware must be placed between [AAuthenticatedMw] and the rest of the request as the user is expected to be present.
type UserInGroupMw struct {
	ggh *query.GetGroupQueryHandler
}

func NewUserInGroupMw(ggh *query.GetGroupQueryHandler) *UserInGroupMw {
	return &UserInGroupMw{ggh}
}

func (u *UserInGroupMw) UserInGroupMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Missing group ID (id) in path parameter", http.StatusBadRequest)
			return
		}
		groupID, err := valueobject.NewGroupIDFromString(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		group, found, err := u.ggh.Handle(r.Context(), query.GetGroupQuery{
			GroupID: groupID,
		})
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, fmt.Sprintf("Group not found for id %s", groupID), http.StatusNotFound)
			return
		}
		if !group.ContainsUser(user.ID) {
			http.Error(w, "Forbidden operation", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
