package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

type JWTHelper struct {
	signatureKey []byte
}

func NewJWTHelper(signatureKey []byte) (*JWTHelper, error) {
	return &JWTHelper{signatureKey}, nil
}

func (h *JWTHelper) NewSignedToken(userID valueobject.UserID) (string, error) {
	// https://golang-jwt.github.io/jwt/usage/create/
	t := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"user_id": userID.String(),
			"exp":     (time.Now().Add(24 * time.Hour)).Unix(),
		})
	token, err := t.SignedString(h.signatureKey)
	return token, err
}

// Verify the token and returns the userID if the token is valid
func (h *JWTHelper) VerifyToken(tokenString string) (*valueobject.UserID, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.signatureKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id not found in token")
	}

	userIDUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	userID, err := valueobject.NewUserID(userIDUUID)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}
