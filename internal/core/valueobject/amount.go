package valueobject

type Amount struct {
	amount   float64
	currency string
}

func NewAmount(amount float64, currency string) (Amount, error) {
	return Amount{amount, currency}, nil
}

func (a Amount) Amount() float64 {
	return a.amount
}

func (a Amount) Currency() string {
	return a.currency
}
