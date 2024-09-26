package auth

import "testing"

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
