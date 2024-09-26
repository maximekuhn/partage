package command

import (
	"context"

	"github.com/maximekuhn/partage/internal/core/contract"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/repositories"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type CreateUser struct {
	Email    valueobject.Email
	Nickname valueobject.Nickname
}

type CreateUserHandler struct {
	uidp     contract.UserIDProvider
	dtp      contract.DatetimeProvider
	userrepo repositories.UserRepo
}

func NewCreateUserHandler(
	uidp contract.UserIDProvider,
	dtp contract.DatetimeProvider,
	userrepo repositories.UserRepo,
) *CreateUserHandler {
	return &CreateUserHandler{uidp, dtp, userrepo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) (valueobject.UserID, error) {
	id := h.uidp.Provide()
	ca := h.dtp.Provide()

	u := entity.NewUser(id, cmd.Email, cmd.Nickname, ca)

	err := h.userrepo.Save(ctx, u)
	if err != nil {
		return id, err
	}

	return id, nil
}
