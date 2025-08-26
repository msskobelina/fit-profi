package access

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRole string

var (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type Token struct {
	jwt.RegisteredClaims
	UserID   int      `json:"userId"`
	UserRole UserRole `json:"userRole"`
}

func EncodeToken(t *Token, hmacSecret string) (string, error) {
	j := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	return j.SignedString([]byte(hmacSecret))
}

func DecodeToken(tokenString string, hmacSecret string) (*Token, error) {
	claims := &Token{}
	parsed, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(tok *jwt.Token) (interface{}, error) {
			if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
			}
			return []byte(hmacSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithLeeway(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("malformed or invalid token: %w", err)
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
