package settlement

type Expense struct {
	amount       uint
	paidBy       string
	participants []string // payer should not be in participants
}

type Debt struct {
	amount uint
	from   string
	to     string
}

// Settle accepts a list of expenses and returns all debts
func Settle(expenses []Expense) []Debt {
	debts := make([]Debt, 0)

	bs := balances(expenses)

	creditors := make(map[string]int, 0) // they are owed money
	debtors := make(map[string]int, 0)   // they owe money

	for p, b := range bs {
		if b > 0 {
			creditors[p] = b
		} else if b < 0 {
			debtors[p] = b
		}
	}

	for len(creditors) > 0 && len(debtors) > 0 {
		cred := largest(creditors, true)
		deb := largest(debtors, false)

		credBalance := creditors[cred]
		debBalance := debtors[deb]

		amount := min(credBalance, -debBalance)

		debts = append(debts, Debt{
			amount: uint(amount),
			from:   deb,
			to:     cred,
		})

		credNewBalance := credBalance - amount
		creditors[cred] = credNewBalance
		if credNewBalance == 0 {
			delete(creditors, cred)
		}

		debNewBalance := debBalance + amount
		debtors[deb] = debNewBalance
		if debNewBalance == 0 {
			delete(debtors, deb)
		}

	}

	return debts
}

func balances(es []Expense) map[string] /* name */ int {
	bs := make(map[string]int)

	for _, e := range es {
		balancedAmount := int(float64(e.amount) / float64(1+len(e.participants)))

		// update balance for the payer
		_, found := bs[e.paidBy]
		if !found {
			bs[e.paidBy] = 0
		}
		bs[e.paidBy] += (int(e.amount) - balancedAmount)

		// update balance for participants
		for _, p := range e.participants {
			_, found = bs[p]
			if !found {
				bs[p] = 0
			}
			bs[p] -= balancedAmount
		}
	}

	return bs
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func largest(m map[string]int, positive bool) string {
	participant := ""
	largest := 0
	for p, b := range m {
		participant = p
		largest = b
		break
	}

	for p, b := range m {
		if positive && b > largest {
			participant = p
			largest = b
		} else if !positive && b < largest {
			participant = p
			largest = b
		}
	}

	return participant
}
