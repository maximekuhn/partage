package web

import (
	"fmt"
	"net/http"

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

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, err := valueobject.NewEmail(r.FormValue("email"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nickname, err := valueobject.NewNickname(r.FormValue("nickname"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password, err := auth.NewPassword(r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := s.authSvc.Hash(password)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cmd := command.CreateUser{
		Email:    email,
		Nickname: nickname,
	}

	ctx := r.Context()

	// TODO: remove internal errors to not expose details to users

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := s.createUserHandler.Handle(ctx, cmd)
	if err != nil {
		_ = tx.Rollback()
		// TODO: error can be user's fault
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.authSvc.Save(ctx, userID, hashedPassword); err != nil {
		_ = tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: should return a page or a redirect
	w.WriteHeader(http.StatusCreated)
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
