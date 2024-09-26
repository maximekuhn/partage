package misc

import (
	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type UserIDProviderProd struct{}

func (p *UserIDProviderProd) Provide() valueobject.UserID {
	id, err := valueobject.NewUserID(uuid.New())
	if err != nil {
		panic(err)
	}
	return id
}
