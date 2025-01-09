package sqlite

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func TestSaveExpense(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteExpenseStore(db)

	tests := []struct {
		title        string
		id           uuid.UUID
		label        string
		paidBy       uuid.UUID
		participants []uuid.UUID
		amount       string
		createdAt    time.Time
		groupID      uuid.UUID
	}{
		{
			title:  "general case",
			id:     uuid.New(),
			label:  "mariokart 8",
			paidBy: uuid.New(),
			participants: []uuid.UUID{
				uuid.New(),
				uuid.New(),
				uuid.New(),
			},
			amount:    "10-$",
			createdAt: time.Now(),
			groupID:   uuid.New(),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			e := createExpense(test.label, test.amount, test.id, test.paidBy, test.groupID, test.participants, test.createdAt)
			if err := s.Save(context.TODO(), e); err != nil {
				t.Errorf("Save(): expected ok got error %v", err)
			}
		})
	}
}

func createExpense(label, amount string,
	id, paidBy, groupID uuid.UUID,
	participants []uuid.UUID,
	createdAt time.Time,
) *entity.Expense {
	l, err := valueobject.NewExpenseLabel(label)
	if err != nil {
		panic(err)
	}

	amountParts := strings.Split(amount, "-")
	if len(amountParts) != 2 {
		panic("expected amount to have format <NUMBER>-<CURRENCY>")
	}
	price, err := strconv.ParseFloat(amountParts[0], 64)
	if err != nil {
		panic(err)
	}
	a, err := valueobject.NewAmount(price, amountParts[1])
	if err != nil {
		panic(err)
	}

	eid, err := valueobject.NewExpenseID(id)
	if err != nil {
		panic(err)
	}

	payerID, err := valueobject.NewUserID(paidBy)
	if err != nil {
		panic(err)
	}

	gid, err := valueobject.NewGroupID(groupID)
	if err != nil {
		panic(err)
	}

	participantsID := make([]valueobject.UserID, 0)
	for _, p := range participants {
		pid, err := valueobject.NewUserID(p)
		if err != nil {
			panic(err)
		}
		participantsID = append(participantsID, pid)
	}

	return entity.NewExpense(eid, l, payerID, participantsID, a, createdAt, gid)
}
