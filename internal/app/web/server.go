package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/maximekuhn/partage/internal/app/web/views"
	"github.com/maximekuhn/partage/internal/user"
)

type Server struct {
	us *user.Service
}

func NewServer(us *user.Service) *Server {
	return &Server{us}
}

func (s *Server) Run() error {
	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/app/web/static"))))
	http.Handle("/", templ.Handler(views.Page("Home")))
	return http.ListenAndServe(":8000", nil)
}
