package valueobject

type ExpenseLabel struct {
	label string
}

func NewExpenseLabel(label string) (ExpenseLabel, error) {
	return ExpenseLabel{label}, nil
}

func (e ExpenseLabel) String() string {
	return e.label
}
