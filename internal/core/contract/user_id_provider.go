package contract

import "github.com/maximekuhn/partage/internal/core/valueobject"

type UserIDProvider interface {
	Provide() valueobject.UserID
}
