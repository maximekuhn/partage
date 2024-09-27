package web

import (
	"net/http"

	"github.com/maximekuhn/partage/internal/app/web/views"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	err := views.Page("Home", views.Index()).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	err := views.Page("Login", views.Login("", msg)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	err := views.Page("Login", views.Register("")).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
