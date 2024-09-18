package settlement

import (
	"fmt"
	"testing"
)

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
					amount: 5,
					from:   "alice",
					to:     "bob",
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
					amount: 10,
					from:   "alice",
					to:     "bob",
				},
				{
					amount: 10,
					from:   "toto",
					to:     "bob",
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
					amount: 4,
					from:   "toto",
					to:     "bob",
				},
				{
					amount: 2,
					from:   "alice",
					to:     "bob",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			debts := Settle(test.expenses)

			got := toMap(debts)
			expected := toMap(test.debts)

			if !compareMaps(expected, got) {
				t.Errorf("Settle(): expected %v got %v", expected, got)
			}
		})
	}
}

// toMap converts a list of debts into a map for easier comparisons.
func toMap(debts []Debt) map[string]uint {
	hs := make(map[string]uint)

	for _, debt := range debts {
		k := fmt.Sprintf("%s -> %s", debt.from, debt.to)
		hs[k] = debt.amount
	}

	return hs
}

func compareMaps(l map[string]uint, r map[string]uint) bool {
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
