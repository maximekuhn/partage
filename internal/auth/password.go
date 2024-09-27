package auth

type Password struct {
	password string
}

func NewPassword(password string) (Password, error) {
	p := Password{password}
	return p, nil
}

func (p Password) String() string {
	return p.password
}

type HashedPassword struct {
	hash []byte // includes salt
}

func NewHashedPassword(hash []byte) HashedPassword {
	return HashedPassword{hash}
}

func (h HashedPassword) Hash() []byte {
	return h.hash
}
