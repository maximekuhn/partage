package misc

import (
	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type GroupIDProviderProd struct{}

func (p *GroupIDProviderProd) Provide() valueobject.GroupID {
	id, err := valueobject.NewGroupID(uuid.New())
	if err != nil {
		panic(err)
	}
	return id
}
