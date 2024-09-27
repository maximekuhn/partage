package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func TestJWTGeneration(t *testing.T) {
	userID, err := valueobject.NewUserID(uuid.New())
	if err != nil {
		panic(err)
	}

	signatureKey := []byte(`64a6988ec0ecacbdf40ecf504e70b9a5f6174a8992c856c7ee22e1e0be03a8890412904b9d17a467d03559fe573c324271615dbcf191e4cfc259b5a01a3bb824`)
	helper, err := NewJWTHelper(signatureKey)
	if err != nil {
		t.Fatalf("NewJWTHelper(): expected ok got error %v", err)
	}

	tokenString, err := helper.NewSignedToken(userID)
	if err != nil {
		t.Fatalf("NewSignedToken(): expected ok got error %v", err)
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		t.Fatal("token has not a valid form")
	}

	claims := parts[1]
	decodedClaims, err := base64.RawURLEncoding.DecodeString(claims)
	if err != nil {
		t.Fatalf("could not decode claims because %v", err)
	}

	type expectedClaims struct {
		User_id string `json:"user_id"`
	}

	var ec expectedClaims

	err = json.NewDecoder(bytes.NewReader(decodedClaims)).Decode(&ec)
	if err != nil {
		t.Fatalf("could not decode claims into JSON because %v", err)
	}

	if ec.User_id != userID.String() {
		t.Fatalf("expected claim user_id to be %s got %s", userID.String(), ec.User_id)
	}
}

func TestValidateToken(t *testing.T) {
	userID, err := valueobject.NewUserID(uuid.New())
	if err != nil {
		panic(err)
	}

	signatureKey := []byte(`64a6988ec0ecacbdf40ecf504e70b9a5f6174a8992c856c7ee22e1e0be03a8890412904b9d17a467d03559fe573c324271615dbcf191e4cfc259b5a01a3bb824`)
	helper, err := NewJWTHelper(signatureKey)
	if err != nil {
		panic(err)
	}

	tokenString, err := helper.NewSignedToken(userID)
	if err != nil {
		panic(err)
	}

	id, err := helper.VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("VerifyToken(): expected ok got error %v", err)
	}

	if *id != userID {
		t.Fatalf("returned user_id: expected %v got %v", userID, *id)
	}
}
