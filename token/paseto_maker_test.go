package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	// PASETO 要求密钥恰好 32 字节
	maker, err := NewPasetoMaker("0123456789abcdef0123456789abcdef")
	require.NoError(t, err)

	username := "alice"
	role := "depositor"
	duration := time.Minute

	token, payload, err := maker.CreateToken(username, role, duration, AccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, username, payload.Username)

	verifiedPayload, err := maker.VerifyToken(token, AccessToken)
	require.NoError(t, err)
	require.Equal(t, payload.ID, verifiedPayload.ID)
	require.Equal(t, username, verifiedPayload.Username)
}

func TestPasetoExpiredToken(t *testing.T) {
	maker, _ := NewPasetoMaker("0123456789abcdef0123456789abcdef")
	token, _, _ := maker.CreateToken("alice", "depositor", -time.Minute, AccessToken)

	_, err := maker.VerifyToken(token, AccessToken)
	require.ErrorIs(t, err, ErrExpiredToken)
}

func TestPasetoInvalidKeyLength(t *testing.T) {
	_, err := NewPasetoMaker("short")
	require.Error(t, err)
}

func TestPasetoWrongTokenType(t *testing.T) {
	maker, _ := NewPasetoMaker("0123456789abcdef0123456789abcdef")
	token, _, _ := maker.CreateToken("alice", "depositor", time.Minute, RefreshToken)

	_, err := maker.VerifyToken(token, AccessToken)
	require.ErrorIs(t, err, ErrInvalidToken)
}
