package store

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type InvitationStore interface {
	Save(ctx context.Context, i valueobject.Invitation) error
	Update(ctx context.Context, i valueobject.Invitation) error
	FindByInviteeID(ctx context.Context, inviteeID valueobject.UserID, groupID valueobject.GroupID) (valueobject.Invitation, bool, error)
}
