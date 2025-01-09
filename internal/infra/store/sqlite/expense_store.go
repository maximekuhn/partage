package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type SQLiteExpenseStore struct {
	db *sql.DB
}

func NewSQLiteExpenseStore(db *sql.DB) *SQLiteExpenseStore {
	return &SQLiteExpenseStore{db}
}

func (s *SQLiteExpenseStore) Save(ctx context.Context, e *entity.Expense) error {
	query := `
    INSERT INTO expense (id, label, payer_id, amount, created_at) VALUES (?, ?, ?, ?, ?)
    `
	dbAmount := fmt.Sprintf("%f#%s", e.Amount.Amount(), e.Amount.Currency())
	if _, err := s.db.ExecContext(ctx, query, e.ID.String(), e.Label.String(), e.PaidBy.String(), dbAmount, e.CreatedAt); err != nil {
		return err
	}

	query = `
    INSERT INTO expense_group (expense_id, group_id) VALUES (?, ?)
    `
	if _, err := s.db.ExecContext(ctx, query, e.ID.String(), e.GroupID.String()); err != nil {
		return err
	}

	query = `
    INSERT INTO expense_user (expense_id, user_id) VALUES (?, ?)
    `
	payerInParticipants := false
	for _, userID := range e.Participants {
		if userID == e.PaidBy {
			payerInParticipants = true
		}
		if _, err := s.db.ExecContext(ctx, query, e.ID.String(), userID.String()); err != nil {
			return err
		}
	}

	if !payerInParticipants {
		if _, err := s.db.ExecContext(ctx, query, e.ID.String(), e.PaidBy.String()); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteExpenseStore) GetAllForGroup(ctx context.Context, groupID valueobject.GroupID) ([]*entity.Expense, error) {
	query := `
    SELECT e.id, e.label, e.payer_id, e.amount, e.created_at, GROUP_CONCAT(eu.user_id)
    FROM expense e
    INNER JOIN expense_group eg ON e.id = eg.expense_id
    INNER JOIN expense_user eu ON e.id = eu.user_id
    WHERE eg.group_id = ?
    GROUP BY e.id
    `
	rows, err := s.db.QueryContext(ctx, query, groupID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*entity.Expense, 0)
	for rows.Next() {
		var expenseID uuid.UUID
		var label string
		var payerID uuid.UUID
		var amount string
		var createdAt time.Time
		var participantsConcat string

		err = rows.Scan(&expenseID, &label, &payerID, &amount, &createdAt,
			&participantsConcat)
		if err != nil {
			return nil, err
		}

		exp, err := tryConvertExpense(expenseID, label, payerID, amount,
			createdAt, participantsConcat, groupID)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, exp)
	}
	return expenses, nil
}

func tryConvertExpense(
	expenseID uuid.UUID,
	label string,
	payerID uuid.UUID,
	amount string,
	createdAt time.Time,
	participantsConcat string,
	groupID valueobject.GroupID,
) (*entity.Expense, error) {
	expID, err := valueobject.NewExpenseID(expenseID)
	if err != nil {
		return nil, err
	}
	expLabel, err := valueobject.NewExpenseLabel(label)
	if err != nil {
		return nil, err
	}
	expPayerID, err := valueobject.NewUserID(payerID)
	if err != nil {
		return nil, err
	}
	expAmount, err := tryConvertAmount(amount)
	if err != nil {
		return nil, err
	}

	participants := strings.Split(participantsConcat, ",")
	expParticipants := make([]valueobject.UserID, 0)
	for _, p := range participants {
		expParticipantID, err := valueobject.NewUserIDFromString(p)
		if err != nil {
			return nil, err
		}
		expParticipants = append(expParticipants, expParticipantID)
	}

	return entity.NewExpense(expID, expLabel, expPayerID, expParticipants,
		expAmount, createdAt, groupID), nil
}

func tryConvertAmount(amount string) (valueobject.Amount, error) {
	a := valueobject.Amount{}
	parts := strings.Split(amount, "#")
	if len(parts) != 2 {
		return a, errors.New("unexpected amount format")
	}
	af64, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return a, err
	}
	return valueobject.NewAmount(af64, parts[1])
}
