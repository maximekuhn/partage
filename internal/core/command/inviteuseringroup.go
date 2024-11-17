package command

import (
	"context"
	"errors"
	"time"

	"github.com/maximekuhn/partage/internal/core/contract"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type InviteUserInGroupCmd struct {
	GroupID      valueobject.GroupID
	UserToInvite valueobject.UserID
}

type InviteUserInGroupCmdHandler struct {
	groupstore      store.GroupStore
	invitationstore store.InvitationStore
	dtprovider      contract.DatetimeProvider
}

func NewInviteUserInGroupCmdHandler(groupstore store.GroupStore, invitationstore store.InvitationStore, dtprovider contract.DatetimeProvider) *InviteUserInGroupCmdHandler {
	return &InviteUserInGroupCmdHandler{groupstore, invitationstore, dtprovider}
}

func (h *InviteUserInGroupCmdHandler) Handle(ctx context.Context, cmd InviteUserInGroupCmd) error {
	g, found, err := h.groupstore.FindByID(ctx, cmd.GroupID)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("group not found")
	}

	for _, member := range g.Members {
		if member == cmd.UserToInvite {
			return errors.New("user is already in the group")
		}
	}

	// `UpdatedAt` field is set to a zero-valued time.Time because it has not been updated yet
	// as the invitation is currently being created
	now := h.dtprovider.Provide()
	i := valueobject.NewInvitation(cmd.UserToInvite, cmd.GroupID, valueobject.InvitationStatusPending, now, time.Time{})
	return h.invitationstore.Save(ctx, i)
}
