package auth

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestBcrypt(t *testing.T) {
	tests := []struct {
		title    string
		password string
	}{
		{
			title:    "simple password",
			password: "admin",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			p, err := NewPassword(test.password)
			if err != nil {
				t.Fatalf("Invalid password: %v", err)
			}

			bcryptHasher := NewBcryptPasswordHasher()

			h, err := bcryptHasher.hash(p)
			if err != nil {
				t.Fatalf("hash(): expected ok got err %v", err)
			}

			verified := bcryptHasher.verify(p, h)
			if !verified {
				t.Fatalf("verify(): expected true got false")
			}
		})
	}

}

func TestSamePasswordDifferentHash(t *testing.T) {
	p, err := NewPassword("Admin1234-secure")
	if err != nil {
		panic(err)
	}

	bcryptHasher := NewBcryptPasswordHasher()

	h1, err := bcryptHasher.hash(p)
	if err != nil {
		t.Fatalf("hash(): expected ok got err %v", err)
	}

	h2, err := bcryptHasher.hash(p)
	if err != nil {
		t.Fatalf("hash(): expected ok got err %v", err)
	}

	if bytes.Equal(h1.hash, h2.hash) {
		t.Fatalf("same password should produce different hashes")
	}
}

func TestHashPasswordTooLong(t *testing.T) {
	// password length > 72 bytes (max length for bcrypt)
	alphabet := []byte("abcdefghijklmnopqrstuvwxyz123456789!@#$%^&*)_+{}[]")

	lengthToGenerate := 100
	tooLongPassword := make([]byte, lengthToGenerate)
	for i := 0; i < lengthToGenerate; i++ {
		tooLongPassword[i] = alphabet[rand.Intn(len(alphabet))]
	}

	password, err := NewPassword(string(tooLongPassword))
	if err != nil {
		panic(err)
	}

	bcryptHasher := NewBcryptPasswordHasher()
	if _, err := bcryptHasher.hash(password); err != nil {
		t.Fatalf("hash(): expected ok got %v", err)
	}
}
