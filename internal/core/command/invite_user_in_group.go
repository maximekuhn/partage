package command

import (
	"context"
	"fmt"

	"github.com/maximekuhn/partage/internal/core/contract"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type InviteUserInGroupCmd struct {
	GroupID      valueobject.GroupID
	UserToInvite valueobject.Email
}

type InviteUserInGroupCmdHandler struct {
	groupstore      store.GroupStore
	userstore       store.UserStore
	invitationstore store.InvitationStore
	dtprovider      contract.DatetimeProvider
}

func NewInviteUserInGroupCmdHandler(groupstore store.GroupStore, userstore store.UserStore, invitationstore store.InvitationStore, dtprovider contract.DatetimeProvider) *InviteUserInGroupCmdHandler {
	return &InviteUserInGroupCmdHandler{groupstore, userstore, invitationstore, dtprovider}
}

func (h *InviteUserInGroupCmdHandler) Handle(ctx context.Context, cmd InviteUserInGroupCmd) error {
	u, found, err := h.userstore.GetByEmail(ctx, cmd.UserToInvite)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("user not found for email %s", cmd.UserToInvite)
	}

	g, found, err := h.groupstore.FindByID(ctx, cmd.GroupID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("group not found for id %s", cmd.GroupID)
	}

	// check if the user is not already in the group
	for _, member := range g.Members {
		if member == u.ID {
			return fmt.Errorf("user with id %s is already in the group", u.ID)
		}
	}

	// create and save pending invitation
	now := h.dtprovider.Provide()
	i := valueobject.NewInvitation(u.ID, g.ID, valueobject.InvitationStatusPending, now, now)
	return h.invitationstore.Save(ctx, i)
}
