package command

import (
	"context"
	"errors"

	"github.com/maximekuhn/partage/internal/core/contract"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type CreateGroupCmd struct {
	Name  valueobject.Groupname
	Owner valueobject.UserID
}

type CreateGroupCmdHandler struct {
	gidp       contract.GroupIDProvider
	dtp        contract.DatetimeProvider
	groupstore store.GroupStore
}

func NewCreateGroupCmdHandler(
	gidp contract.GroupIDProvider,
	dtp contract.DatetimeProvider,
	groupstore store.GroupStore,
) *CreateGroupCmdHandler {
	return &CreateGroupCmdHandler{gidp, dtp, groupstore}
}

func (h *CreateGroupCmdHandler) Handle(ctx context.Context, cmd CreateGroupCmd) (valueobject.GroupID, error) {
	id := h.gidp.Provide()

	_, found, err := h.groupstore.FindByName(ctx, cmd.Name)
	if err != nil {
		return id, err
	}
	if found {
		return id, errors.New("another group with the same name already exists")
	}

	ca := h.dtp.Provide()
	g := entity.NewGroup(id, cmd.Name, make([]valueobject.UserID, 0), cmd.Owner, ca)

	err = h.groupstore.Save(ctx, g)

	return id, err
}
