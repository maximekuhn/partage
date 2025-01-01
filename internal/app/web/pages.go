package web

import (
	"fmt"
	"net/http"

	"github.com/maximekuhn/partage/internal/app/web/middleware"
	"github.com/maximekuhn/partage/internal/app/web/views"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/query"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	authmwdata, ok := r.Context().Value(middleware.AuthDatacontextKey).(middleware.AuthMwData)

	var user *entity.User
	if ok && authmwdata.Authenticated {
		user = authmwdata.User
	}

	err := views.Page("Home", user, views.Index(user)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	authmwdata, ok := r.Context().Value(middleware.AuthDatacontextKey).(middleware.AuthMwData)

	var user *entity.User
	if ok && authmwdata.Authenticated {
		user = authmwdata.User
	}

	msg := r.URL.Query().Get("msg")
	errMsg := r.URL.Query().Get("err_msg")
	err := views.Page("Login", user, views.Login(errMsg, msg)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	authmwdata, ok := r.Context().Value(middleware.AuthDatacontextKey).(middleware.AuthMwData)

	var user *entity.User
	if ok && authmwdata.Authenticated {
		user = authmwdata.User
	}
	err := views.Page("Register", user, views.Register("")).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGroups(w http.ResponseWriter, r *http.Request) {
	authmwdata, ok := r.Context().Value(middleware.AuthDatacontextKey).(middleware.AuthMwData)

	var user *entity.User
	if ok && authmwdata.Authenticated {
		user = authmwdata.User
	}
	if user == nil {
		panic("user is null, middleware should have returned earlier")
	}

	groups, err := s.app.GetGroupsHandler.Handle(r.Context(), query.GetGroupsForUserQuery{
		UserID: user.ID,
	})
	if err != nil {
		// TODO: return to a common page '500 Internal Server Error'
		groups = make([]entity.Group, 0)
	}

	err = views.Page("Groups", user, views.Groups(user, groups, "")).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGroup(w http.ResponseWriter, r *http.Request) {
	authmwdata, ok := r.Context().Value(middleware.AuthDatacontextKey).(middleware.AuthMwData)

	var user *entity.User
	if ok && authmwdata.Authenticated {
		user = authmwdata.User
	}
	if user == nil {
		panic("user is null, middleware should have returned earlier")
	}
	id := r.PathValue("id")
	groupID, err := valueobject.NewGroupIDFromString(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	g, found, err := s.app.GetGroupDetailsHandler.Handle(r.Context(), query.GetGroupDetailsQuery{
		GroupID: groupID,
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	err = views.Page("Group", user, views.Group(g)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
