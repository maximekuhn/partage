package command

import (
	"context"
	"errors"

	"github.com/maximekuhn/partage/internal/core/contract"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/store"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type AddExpenseToGroupCmd struct {
	GroupID valueobject.GroupID
	PayerID valueobject.UserID
	Label   valueobject.ExpenseLabel
	Amount  valueobject.Amount

	// AllMembers is set to true if all group members are
	// participating to this expense.
	AllMembers bool

	// Participants is only used when AllMembers is set to false.
	Participants []valueobject.UserID
}

type AddExpenseToGroupCmdHandler struct {
	groupstore        store.GroupStore
	expensestore      store.ExpenseStore
	expenseidprovider contract.ExpenseIDProvider
	dtprovider        contract.DatetimeProvider
}

func NewAddExpenseToGroupCmdHandler(groupstore store.GroupStore,
	expensestore store.ExpenseStore,
	expenseidprovider contract.ExpenseIDProvider,
	dtprovider contract.DatetimeProvider) *AddExpenseToGroupCmdHandler {
	return &AddExpenseToGroupCmdHandler{groupstore, expensestore, expenseidprovider, dtprovider}
}

func (h *AddExpenseToGroupCmdHandler) Handle(ctx context.Context, cmd AddExpenseToGroupCmd) error {
	g, found, err := h.groupstore.FindByID(ctx, cmd.GroupID)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("group not found")
	}

	var participants []valueobject.UserID
	if cmd.AllMembers {
		participants = g.Members
	} else {
		participants = make([]valueobject.UserID, 0)
		participants = append(participants, cmd.Participants...)
	}

	expenseid := h.expenseidprovider.Provide()
	createdAt := h.dtprovider.Provide()
	e := entity.NewExpense(expenseid, cmd.Label, cmd.PayerID, participants, cmd.Amount, createdAt, cmd.GroupID)

	if err := h.expensestore.Save(ctx, e); err != nil {
		return err
	}

	return nil
}
