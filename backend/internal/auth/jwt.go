package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID int64 `json:"uid"`
	Role   Role  `json:"role"`
	Banned bool  `json:"banned"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) *JWT {
	return &JWT{secret: []byte(secret), ttl: ttl}
}

func (j *JWT) Mint(userID int64, role Role, banned bool) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		UserID: userID,
		Role:   role,
		Banned: banned,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

func (j *JWT) Parse(tokenStr string) (int64, Role, bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return j.secret, nil
	})
	if err != nil {
		return 0, "", false, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", false, ErrInvalidToken
	}

	return claims.UserID, claims.Role, claims.Banned, nil
}
