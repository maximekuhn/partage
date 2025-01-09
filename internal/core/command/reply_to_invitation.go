package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/maximekuhn/partage/internal/core/contract"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type ReplyToInvitationCmd struct {
	Invitee  valueobject.UserID
	GroupID  valueobject.GroupID
	Accepted bool
}

type ReplyToInvitationCmdHandler struct {
	invitationstore store.InvitationStore
	groupstore      store.GroupStore
	dtprovider      contract.DatetimeProvider
}

func NewReplyToInvitationCmdHandler(
	invitationstore store.InvitationStore,
	dtprovider contract.DatetimeProvider,
	groupstore store.GroupStore,
) *ReplyToInvitationCmdHandler {
	return &ReplyToInvitationCmdHandler{invitationstore, groupstore, dtprovider}
}

func (h *ReplyToInvitationCmdHandler) Handle(ctx context.Context, cmd ReplyToInvitationCmd) error {
	ivt, found, err := h.invitationstore.FindByInviteeID(ctx, cmd.Invitee, cmd.GroupID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("Invitation not found for user")
	}

	// invitation must be in pending status
	if !ivt.IsPending() {
		return errors.New("Invitation is not in pending status")
	}

	if cmd.Accepted {
		return h.acceptInvitation(ctx, ivt)
	}
	return h.rejectInvitation(ctx, ivt)
}

func (h *ReplyToInvitationCmdHandler) acceptInvitation(
	ctx context.Context,
	ivt valueobject.Invitation,
) error {
	grp, found, err := h.groupstore.FindByID(ctx, ivt.GroupID)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("group not found")
	}
	grp.Members = append(grp.Members, ivt.UserID)
	if err := h.groupstore.Update(ctx, grp); err != nil {
		return err
	}
	return h.updateInvitation(ctx, ivt, true)
}

func (h *ReplyToInvitationCmdHandler) rejectInvitation(
	ctx context.Context,
	ivt valueobject.Invitation,
) error {
	return h.updateInvitation(ctx, ivt, false)
}

func (h *ReplyToInvitationCmdHandler) updateInvitation(
	ctx context.Context,
	ivt valueobject.Invitation,
	accepted bool,
) error {
	if accepted {
		ivt.Status = valueobject.InvitationStatusAccepted
	} else {
		ivt.Status = valueobject.InvitationStatusRejected
	}

	updatedAt := h.dtprovider.Provide()
	ivt.UpdatedAt = updatedAt

	return h.invitationstore.Update(ctx, ivt)
}
