package auth

import (
	"bytes"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

func (b *BcryptPasswordHasher) hash(p Password) (HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.password), bcrypt.DefaultCost)
	h := HashedPassword{hash}
	if err != nil {
		return h, err
	}
	return h, nil
}

// verify if the password is the same as the hash
func (b *BcryptPasswordHasher) verify(p Password, h HashedPassword) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}
	if bytes.Equal(hash, h.hash) {
		return false
	}
	return true
}