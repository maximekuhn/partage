package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maximekuhn/partage/internal/core/entity"
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
