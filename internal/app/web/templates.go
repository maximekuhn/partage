package web

import (
	"net/http"

	"github.com/maximekuhn/partage/internal/app/web/middleware"
	"github.com/maximekuhn/partage/internal/app/web/views"
	"github.com/maximekuhn/partage/internal/core/entity"
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

	err := views.Page("Groups", user, views.Groups(user, "")).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

}
