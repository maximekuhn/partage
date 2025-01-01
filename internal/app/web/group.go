package web

import (
	"fmt"
	"net/http"

	"github.com/maximekuhn/partage/internal/app/web/middleware"
	"github.com/maximekuhn/partage/internal/app/web/views"
	"github.com/maximekuhn/partage/internal/core/command"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/query"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func (s *Server) handleCreateGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s /group/create\n", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	authmwdata, ok := r.Context().Value(middleware.AuthDatacontextKey).(middleware.AuthMwData)

	var user *entity.User
	if ok && authmwdata.Authenticated {
		user = authmwdata.User
	}

	w.Header().Add("Content-Type", "text/html")
	ctx := r.Context()

	if user == nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	groups, err := s.app.GetGroupsHandler.Handle(r.Context(), query.GetGroupsForUserQuery{
		UserID: user.ID,
	})

	if err != nil {
		// TODO: return to a common page '500 Internal Server Error'
		groups = make([]entity.Group, 0)
	}

	if err := r.ParseForm(); err != nil {
		_ = views.Page("Groups", user, views.Groups(user, groups, "")).Render(ctx, w)
		return
	}

	groupname, err := valueobject.NewGroupname(r.FormValue("group_name"))
	if err != nil {
		_ = views.Page("Groups", user, views.Groups(user, groups, err.Error())).Render(ctx, w)
		return
	}

	tx, err := s.app.BeginTx(ctx)
	if err != nil {
		_ = views.Page("Groups", user, views.Groups(user, groups, "Something went wrong")).Render(ctx, w)
		return
	}

	_, err = s.app.CreateGroupHandler.Handle(ctx, command.CreateGroupCmd{
		Name:  groupname,
		Owner: user.ID,
	})
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("failed to create group: %s\n", err)
		_ = views.Page("Groups", user, views.Groups(user, groups, "Oops! Something went wrong and your group could not be created :(")).Render(ctx, w)
		return
	}

	if err := tx.Commit(); err != nil {
		_ = views.Page("Groups", user, views.Groups(user, groups, "Something went wrong")).Render(ctx, w)
		return
	}

	http.Redirect(w, r, "/group", http.StatusSeeOther)
}
