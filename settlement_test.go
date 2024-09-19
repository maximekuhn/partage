package settlement

import (
	"fmt"
	"slices"
	"testing"
)

func TestExpenseCtor(t *testing.T) {
	tests := []struct {
		title        string
		amount       uint
		paidBy       string
		participants []string
		shouldErr    bool
	}{
		{
			title:        "ok",
			amount:       10,
			paidBy:       "alice",
			participants: []string{"bob", "toto"},
			shouldErr:    false,
		},
		{
			title:        "err payer in participants",
			amount:       10,
			paidBy:       "alice",
			participants: []string{"alice", "bob", "toto"},
			shouldErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			e, err := NewExpense(test.amount, test.paidBy, test.participants)

			if test.shouldErr && err != nil {
				return
			}

			if !test.shouldErr && err != nil {
				t.Errorf("NewExpense(): expected no error got %v", err)
			}

			if test.shouldErr && err == nil {
				t.Error("NewExpense(): expected an error got nothing")
			}

			if e.amount != test.amount {
				t.Errorf("NewExpense(): expected amount %v got %v", test.amount, e.amount)
			}

			if e.paidBy != test.paidBy {
				t.Errorf("NewExpense(): expected paidBy %v got %v", test.paidBy, e.paidBy)
			}

			if !slicesEqual(e.participants, test.participants) {
				t.Errorf(
					"NewExpense(): expected participants %v got %v",
					test.participants, e.participants,
				)
			}

		})
	}
}

func TestSettle(t *testing.T) {
	tests := []struct {
		title    string
		expenses []Expense
		debts    []Debt // expected
	}{
		// --
		{
			title:    "no expense no debt",
			expenses: []Expense{},
			debts:    []Debt{},
		},
		// --
		{
			title: "alice owes bob 5",
			expenses: []Expense{
				{
					amount:       10,
					paidBy:       "bob",
					participants: []string{"alice"},
				},
			},
			debts: []Debt{
				{
					Amount: 5,
					From:   "alice",
					To:     "bob",
				},
			},
		},
		// --
		{
			title: "already settled",
			expenses: []Expense{
				{
					amount:       10,
					paidBy:       "bob",
					participants: []string{"alice"},
				},
				{
					amount:       10,
					paidBy:       "alice",
					participants: []string{"bob"},
				},
			},
			debts: []Debt{},
		},
		// --
		{
			title: "3 persons",
			expenses: []Expense{
				{
					amount:       30,
					paidBy:       "bob",
					participants: []string{"alice", "toto"},
				},
			},
			debts: []Debt{
				{
					Amount: 10,
					From:   "alice",
					To:     "bob",
				},
				{
					Amount: 10,
					From:   "toto",
					To:     "bob",
				},
			},
		},
		// --
		{
			title: "3 persons - multiple expenses",
			expenses: []Expense{
				{
					amount:       12,
					paidBy:       "bob",
					participants: []string{"alice", "toto"},
				},
				{
					amount:       4,
					paidBy:       "alice",
					participants: []string{"bob"},
				},
			},
			debts: []Debt{
				{
					Amount: 4,
					From:   "toto",
					To:     "bob",
				},
				{
					Amount: 2,
					From:   "alice",
					To:     "bob",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			debts := Settle(test.expenses)

			got := toMap(debts)
			expected := toMap(test.debts)

			if !mapsEqual(expected, got) {
				t.Errorf("Settle(): expected %v got %v", expected, got)
			}
		})
	}
}

// toMap converts a list of debts into a map for easier comparisons.
func toMap(debts []Debt) map[string]uint {
	hs := make(map[string]uint)

	for _, debt := range debts {
		k := fmt.Sprintf("%s -> %s", debt.From, debt.To)
		hs[k] = debt.Amount
	}

	return hs
}

func mapsEqual(l map[string]uint, r map[string]uint) bool {
	if len(l) != len(r) {
		return false
	}

	for kl, vl := range l {
		vr, found := r[kl]
		if !found {
			return false
		}
		if vr != vl {
			return false
		}
	}

	return true
}

func slicesEqual(l []string, r []string) bool {
	if len(l) != len(r) {
		return false
	}

	slices.Sort(l)
	slices.Sort(r)

	for i, v := range l {
		if r[i] != v {
			return false
		}
	}

	return true
}
