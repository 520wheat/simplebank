package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

// jwtClaims 是 JWT 格式的 claims，实现 jwt.Claims 接口
type jwtClaims struct {
	ID        string    `json:"jti"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

func (maker *JWTMaker) toClaims(payload *Payload) *jwtClaims {
	return &jwtClaims{
		ID:        payload.ID.String(),
		Username:  payload.Username,
		Role:      payload.Role,
		TokenType: payload.TokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
		},
	}
}

func (claims *jwtClaims) toPayload(expectedType TokenType) (*Payload, error) {
	if claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        id,
		Username:  claims.Username,
		Role:      claims.Role,
		IssuedAt:  claims.RegisteredClaims.IssuedAt.Time,
		ExpiredAt: claims.RegisteredClaims.ExpiresAt.Time,
		TokenType: claims.TokenType,
	}, nil
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters, got %d",
			minSecretKeySize, len(secretKey))
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration, tokenType)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, maker.toClaims(payload))
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(tokenStr string, tokenType TokenType) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*jwtClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims.toPayload(tokenType)
}
