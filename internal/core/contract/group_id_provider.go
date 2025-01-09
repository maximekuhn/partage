package contract

import "github.com/maximekuhn/partage/internal/core/valueobject"

type GroupIDProvider interface {
	Provide() valueobject.GroupID
}
