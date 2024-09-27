package web

import (
	"fmt"
	"net/http"

	"github.com/maximekuhn/partage/internal/app/web/views"
	"github.com/maximekuhn/partage/internal/auth"
	"github.com/maximekuhn/partage/internal/core/command"
	"github.com/maximekuhn/partage/internal/core/query"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func (s *Server) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s /register\n", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		views.Page("Register", views.Register("Some informations are missing !")).Render(ctx, w)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, err := valueobject.NewEmail(r.FormValue("email"))
	if err != nil {
		views.Register("Please enter a valid email").Render(ctx, w)
		return
	}

	nickname, err := valueobject.NewNickname(r.FormValue("nickname"))
	if err != nil {
		views.Page("Register", views.Register("Please enter a valid nickname")).Render(ctx, w)
		return
	}

	password, err := auth.NewPassword(r.FormValue("password"))
	if err != nil {
		views.Page("Register", views.Register("Password is not strong enough")).Render(ctx, w)
		return
	}
	passwordConfirm, err := auth.NewPassword(r.FormValue("password_confirm"))
	if err != nil {
		views.Page("Register", views.Register("Password is not strong enough")).Render(ctx, w)
		return
	}
	if password != passwordConfirm {
		views.Page("Register", views.Register("Password and confirmation don't match")).Render(ctx, w)
		return
	}
	hashedPassword, err := s.authSvc.Hash(password)
	if err != nil {
		views.Page("Register", views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	cmd := command.CreateUser{
		Email:    email,
		Nickname: nickname,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		views.Page("Register", views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	userID, err := s.createUserHandler.Handle(ctx, cmd)
	if err != nil {
		_ = tx.Rollback()
		// TODO: error can be user's fault
		views.Page("Register", views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	if err := s.authSvc.Save(ctx, userID, hashedPassword); err != nil {
		_ = tx.Rollback()
		// XXX: can it be user's fault here too?
		views.Page("Register", views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	if err = tx.Commit(); err != nil {
		views.Page("Register", views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s /login\n", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, err := valueobject.NewEmail(r.FormValue("email"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password, err := auth.NewPassword(r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, found, err := s.getUserByEmailHandler.Handle(ctx, query.GetUserByEmailCommand{
		Email: email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	authenticated := s.authSvc.Authenticate(ctx, u.ID, password)

	if !authenticated {
		http.Error(w, "you are not who you pretend to be or you don't exist", http.StatusForbidden)
		return
	}

	jwt, err := s.authSvc.GenerateJWT(u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "jwt: %s", jwt)
}
