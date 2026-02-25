package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) (*Manager, error) {
	if secret == "" {
		return nil, errors.New("JWT_SECRET is empty")
	}
	return &Manager{secret: []byte(secret), ttl: ttl}, nil
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (m *Manager) NewToken(userID string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

func (m *Manager) ParseToken(tokenString string) (string, error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return "", errors.New("invalid token")
	}
	if claims.UserID == "" {
		return "", errors.New("empty user_id in token")
	}
	return claims.UserID, nil
}
