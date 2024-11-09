package auth

import (
	"crypto/sha512"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

func (b *BcryptPasswordHasher) hash(p Password) (HashedPassword, error) {
	passwordBytes := []byte(p.password)
	if len(passwordBytes) >= 72 {
		hashedPasswordBytes := sha512.New().Sum(passwordBytes)
		passwordBytes = hashedPasswordBytes[len(passwordBytes):]
	}

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	h := HashedPassword{hash}
	if err != nil {
		return h, err
	}
	return h, nil
}

// verify if the password is the same as the hash
func (b *BcryptPasswordHasher) verify(p Password, h HashedPassword) bool {
	return bcrypt.CompareHashAndPassword(h.hash, []byte(p.password)) == nil
}
