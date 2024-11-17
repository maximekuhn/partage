package store

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type InvitationStore interface {
	Save(ctx context.Context, i valueobject.Invitation) error
}
