package web

import (
	"context"
	"database/sql"

	"github.com/maximekuhn/partage/internal/auth"
	"github.com/maximekuhn/partage/internal/core/command"
	"github.com/maximekuhn/partage/internal/core/query"
	"github.com/maximekuhn/partage/internal/infra/misc"
	"github.com/maximekuhn/partage/internal/infra/store/sqlite"
)

type application struct {
	db                     *sql.DB
	AuthService            *auth.AuthService
	CreateUserHandler      *command.CreateUserHandler
	GetUserByEmailHandler  *query.GetUserByEmailQueryHandler
	GetUserByIDHandler     *query.GetUserByIDQueryHandler
	CreateGroupHandler     *command.CreateGroupCmdHandler
	GetGroupsHandler       *query.GetGroupsForUserQueryHandler
	GetGroupHandler        *query.GetGroupQueryHandler
	GetGroupDetailsHandler *query.GetGroupDetailsQueryHandler
}

func newApplication(dbFilepath, jwtSignatureKey string) (*application, error) {
	// db stuff
	db, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err = sqlite.ApplyMigrations(db); err != nil {
		return nil, err
	}

	// auth service
	jwtHelper, err := auth.NewJWTHelper([]byte(jwtSignatureKey))
	if err != nil {
		return nil, err
	}

	authSvc := auth.NewAuthService(
		auth.NewBcryptPasswordHasher(),
		sqlite.NewSQLiteAuthStore(db),
		jwtHelper,
	)

	// stores
	userstore := sqlite.NewSQLiteUserStore(db)
	groupstore := sqlite.NewSQLiteGroupStore(db)
	expensestore := sqlite.NewSQLiteExpenseStore(db)

	// commands handlers
	createUserHandler := command.NewCreateUserHandler(
		&misc.UserIDProviderProd{},
		&misc.DatetimeProviderProd{},
		userstore)

	createGroupHandler := command.NewCreateGroupCmdHandler(
		&misc.GroupIDProviderProd{},
		&misc.DatetimeProviderProd{},
		groupstore)

	// queries handlers
	getUserByEmailHandler := query.NewGetUserByEmailCommandHandler(userstore)
	getUserByIDHandler := query.NewGetUserByIDCommandHandler(userstore)
	getGroupsHandler := query.NewGetGroupsForUserQueryHandler(groupstore)
	getGroupHandler := query.NewGetGroupQueryHandler(groupstore)
	getGroupDetailsHandler := query.NewGetGroupDetailsQueryHandler(groupstore, userstore, expensestore)

	return &application{
		db,
		authSvc,
		createUserHandler,
		getUserByEmailHandler,
		getUserByIDHandler,
		createGroupHandler,
		getGroupsHandler,
		getGroupHandler,
		getGroupDetailsHandler,
	}, nil
}

func (a *application) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return a.db.BeginTx(ctx, nil)
}
