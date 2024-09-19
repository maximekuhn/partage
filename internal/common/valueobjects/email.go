package valueobjects

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	e := Email{email}
	return e, nil
}

func (e Email) String() string {
	return e.email
}
