package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// JWTAuth はJWTベースの認証を扱います
type JWTAuth struct {
	secret []byte
}

// jwtClaims はカスタムクレームです
type jwtClaims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

// NewJWTAuth は新しいJWTAuthを作成します
func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{
		secret: []byte(secret),
	}
}

// GenerateToken はユーザーIDからJWTトークンを生成します
func (j *JWTAuth) GenerateToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &jwtClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}

// ParseToken はJWTトークンを検証し、ユーザーIDを返します
func (j *JWTAuth) ParseToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id in token: %w", err)
	}

	return userID, nil
}


