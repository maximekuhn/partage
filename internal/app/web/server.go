package web

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/partage/internal/app/web/views"
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

	return &Server{db, authSvc, createUserHandler, getUserByEmailHandler}, nil
}

func (s *Server) Run() error {
	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/app/web/static"))))

	http.Handle("/", templ.Handler(views.Page("Home", views.Index())))
	http.Handle("GET /register", templ.Handler(views.Page("Register", views.Register(""))))
	http.HandleFunc("POST /register", s.handleRegisterUser)
	http.Handle("GET /login", templ.Handler(views.Page("Login", views.Login())))
	http.HandleFunc("POST /login", s.handleLoginUser)

	fmt.Println("server is up and running")

	return http.ListenAndServe(":8000", nil)
}
