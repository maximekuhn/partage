package web

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/partage/internal/app/web/middleware"
	"github.com/maximekuhn/partage/internal/auth"
	"github.com/maximekuhn/partage/internal/core/command"
	"github.com/maximekuhn/partage/internal/core/query"
	"github.com/maximekuhn/partage/internal/infra/misc"
	"github.com/maximekuhn/partage/internal/infra/store/sqlite"
)

type Server struct {
	db                    *sql.DB
	authSvc               *auth.AuthService
	createUserHandler     *command.CreateUserHandler
	getUserByEmailHandler *query.GetUserByEmailCommandHandler
	getUserByIDHandler    *query.GetUserByIDCommandHandler
}

func NewServer(config ServerConfig) (*Server, error) {
	fmt.Printf("config: %v\n", config)
	db, err := sql.Open("sqlite3", config.DBFilepath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err = sqlite.ApplyMigrations(db); err != nil {
		return nil, err
	}

	jwtHelper, err := auth.NewJWTHelper(config.JWTSignatureKey)
	if err != nil {
		return nil, err
	}

	authSvc := auth.NewAuthService(
		auth.NewBcryptPasswordHasher(),
		sqlite.NewSQLiteAuthStore(db),
		jwtHelper,
	)

	userstore := sqlite.NewSQLiteUserStore(db)

	createUserHandler := command.NewCreateUserHandler(
		&misc.UserIDProviderProd{},
		&misc.DatetimeProviderProd{},
		userstore,
	)

	getUserByEmailHandler := query.NewGetUserByEmailCommandHandler(userstore)
	getUserByIDHandler := query.NewGetUserByIDCommandHandler(userstore)

	return &Server{db, authSvc, createUserHandler, getUserByEmailHandler, getUserByIDHandler}, nil
}

func (s *Server) Run() error {
	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/app/web/static"))))

	authMw := middleware.NewAuthMw(s.authSvc, s.getUserByIDHandler)

	http.Handle("/", authMw.AuthMiddleware(http.HandlerFunc(s.handleIndex)))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("GET /register", s.handleRegister)
	http.HandleFunc("POST /register", s.handleRegisterUser)
	http.HandleFunc("GET /login", s.handleLogin)
	http.HandleFunc("POST /login", s.handleLoginUser)
	http.Handle("POST /logout", authMw.AuthMiddleware(http.HandlerFunc(s.handleLogoutUser)))

	fmt.Println("server is up and running")

	return http.ListenAndServe(":8000", nil)
}
