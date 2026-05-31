package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

// Payload 就是 JWT/PASETO 的 payload（claims）部分
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	TokenType TokenType `json:"token_type"`
}

// NewPayload 创建一个新的 Token payload
func NewPayload(username string, role string, duration time.Duration, tokenType TokenType) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		TokenType: tokenType,
	}
	return payload, nil
}

// Valid 检查 token 是否过期
func (payload *Payload) Valid(expectedType TokenType) error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	if payload.TokenType != expectedType {
		return ErrInvalidToken
	}
	return nil
}
