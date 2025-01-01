package web

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/maximekuhn/partage/internal/app/web/middleware"
	"github.com/maximekuhn/partage/internal/app/web/views"
	"github.com/maximekuhn/partage/internal/auth"
	"github.com/maximekuhn/partage/internal/core/command"
	"github.com/maximekuhn/partage/internal/core/entity"
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
		_ = views.Page("Register", nil, views.Register("Some informations are missing !")).Render(ctx, w)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, err := valueobject.NewEmail(r.FormValue("email"))
	if err != nil {
		_ = views.Register("Please enter a valid email").Render(ctx, w)
		return
	}

	nickname, err := valueobject.NewNickname(r.FormValue("nickname"))
	if err != nil {
		_ = views.Page("Register", nil, views.Register("Please enter a valid nickname")).Render(ctx, w)
		return
	}

	password, err := auth.NewPassword(r.FormValue("password"))
	if err != nil {
		_ = views.Page("Register", nil, views.Register("Password is not strong enough")).Render(ctx, w)
		return
	}
	passwordConfirm, err := auth.NewPassword(r.FormValue("confirm_password"))
	if err != nil {
		_ = views.Page("Register", nil, views.Register("Password is not strong enough")).Render(ctx, w)
		return
	}
	if password != passwordConfirm {
		_ = views.Page("Register", nil, views.Register("Password and confirmation don't match")).Render(ctx, w)
		return
	}
	hashedPassword, err := s.app.AuthService.Hash(password)
	if err != nil {
		_ = views.Page("Register", nil, views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	cmd := command.CreateUserCmd{
		Email:    email,
		Nickname: nickname,
	}

	tx, err := s.app.BeginTx(ctx)
	if err != nil {
		_ = views.Page("Register", nil, views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	userID, err := s.app.CreateUserHandler.Handle(ctx, cmd)
	if err != nil {
		_ = tx.Rollback()
		// TODO: error can be user's fault
		_ = views.Page("Register", nil, views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	if err := s.app.AuthService.Save(ctx, userID, hashedPassword); err != nil {
		_ = tx.Rollback()
		// XXX: can it be user's fault here too?
		_ = views.Page("Register", nil, views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	if err = tx.Commit(); err != nil {
		_ = views.Page("Register", nil, views.Register("Something went wrong :( Please try again")).Render(ctx, w)
		return
	}

	msg := url.QueryEscape("Account created successfully !")
	redirectUrl := fmt.Sprintf("/login?msg=%s", msg)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (s *Server) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s /login\n", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, err := valueobject.NewEmail(r.FormValue("email"))
	if err != nil {
		_ = views.Page("Login", nil, views.Login("Please enter a valid email", "")).Render(ctx, w)
		return
	}

	password, err := auth.NewPassword(r.FormValue("password"))
	if err != nil {
		_ = views.Page("Login", nil, views.Login("Invalid password (not strong enough)", "")).Render(ctx, w)
		return
	}

	u, found, err := s.app.GetUserByEmailHandler.Handle(ctx, query.GetUserByEmailQuery{
		Email: email,
	})
	if err != nil {
		_ = views.Page("Login", nil, views.Login("Something went wrong :( Please try again later.", "")).Render(ctx, w)
		return
	}
	if !found {
		_ = views.Page("Login", nil, views.Login("Invalid credentials or account not found", "")).Render(ctx, w)
		return
	}

	authenticated := s.app.AuthService.Authenticate(ctx, u.ID, password)

	if !authenticated {
		_ = views.Page("Login", nil, views.Login("Invalid credentials or account not found", "")).Render(ctx, w)
		return
	}

	jwt, err := s.app.AuthService.GenerateJWT(u.ID)
	if err != nil {
		_ = views.Page("Login", nil, views.Login("Something went wrong :( Please try again later.", "")).Render(ctx, w)
		return
	}

	// Store JWT in cookies
	// is it ok?
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    jwt,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s /logout\n", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	authmwdata, ok := r.Context().Value(middleware.AuthDatacontextKey).(middleware.AuthMwData)

	var user *entity.User
	if ok && authmwdata.Authenticated {
		user = authmwdata.User
	}

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Remove JWT from cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
