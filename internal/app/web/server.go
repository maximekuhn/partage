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
	getUserByEmailHandler *query.GetUserByEmailQueryHandler
	getUserByIDHandler    *query.GetUserByIDQueryHandler
	createGroupHandler    *command.CreateGroupCmdHandler
	getGroupsHandler      *query.GetGroupsForUserQueryHandler
	getGroupHandler       *query.GetGroupQueryHandler
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

	groupstore := sqlite.NewSQLiteGroupStore(db)

	createGroupHandler := command.NewCreateGroupCmdHandler(
		&misc.GroupIDProviderProd{},
		&misc.DatetimeProviderProd{},
		groupstore,
	)

	getGroupsHandler := query.NewGetGroupsForUserQueryHandler(groupstore)
	getGroupHandler := query.NewGetGroupQueryHandler(groupstore)

	return &Server{
		db,
		authSvc,
		createUserHandler,
		getUserByEmailHandler,
		getUserByIDHandler,
		createGroupHandler,
		getGroupsHandler,
		getGroupHandler,
	}, nil
}

func (s *Server) Run() error {
	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/app/web/static"))))
	http.HandleFunc("/favicon.ico", faviconHandler)

	authMw := middleware.NewAuthMw(s.authSvc, s.getUserByIDHandler)
	authenticatedMw := middleware.NewAuthenticatedMw()
	userInGroupMw := middleware.NewUserInGroupMw(s.getGroupHandler)

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
