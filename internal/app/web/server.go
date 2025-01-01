package web

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/partage/internal/app/web/middleware"
)

type Server struct {
	app *application
}

func NewServer(config ServerConfig) (*Server, error) {
	app, err := newApplication(config.DBFilepath, string(config.JWTSignatureKey))
	if err != nil {
		return nil, err
	}
	return &Server{app}, nil
}

func (s *Server) Run() error {
	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/app/web/static"))))
	http.HandleFunc("/favicon.ico", faviconHandler)

	authMw := middleware.NewAuthMw(s.app.AuthService, s.app.GetUserByIDHandler)
	authenticatedMw := middleware.NewAuthenticatedMw()
	userInGroupMw := middleware.NewUserInGroupMw(s.app.GetGroupHandler)

	http.Handle("/", authMw.AuthMiddleware(s.handleIndex))
	http.HandleFunc("GET /register", s.handleRegister)
	http.HandleFunc("POST /register", s.handleRegisterUser)
	http.HandleFunc("GET /login", s.handleLogin)
	http.HandleFunc("POST /login", s.handleLoginUser)
	http.Handle("POST /logout", authMw.AuthMiddleware(s.handleLogoutUser))
	http.Handle("GET /group", authMw.AuthMiddleware(authenticatedMw.AuthenticatedMiddleware(s.handleGroups)))
	http.Handle("GET /group/{id}", authMw.AuthMiddleware(authenticatedMw.AuthenticatedMiddleware(userInGroupMw.UserInGroupMiddleware(s.handleGroup))))
	http.Handle("POST /group/create", authMw.AuthMiddleware(authenticatedMw.AuthenticatedMiddleware(s.handleCreateGroup)))

	fmt.Println("server is up and running")

	return http.ListenAndServe(":8000", nil)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /favicon.ico")
	http.ServeFile(w, r, "./internal/app/web/static/favicon.ico")
}
