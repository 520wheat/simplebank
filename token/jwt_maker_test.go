package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker("0123456789abcdef0123456789abcdef")
	require.NoError(t, err)

	username := "alice"
	role := "depositor"
	duration := time.Minute

	token, payload, err := maker.CreateToken(username, role, duration, AccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.Equal(t, AccessToken, payload.TokenType)

	// 验证 token
	verifiedPayload, err := maker.VerifyToken(token, AccessToken)
	require.NoError(t, err)
	require.Equal(t, payload.ID, verifiedPayload.ID)
	require.Equal(t, username, verifiedPayload.Username)
}

func TestJWTExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker("0123456789abcdef0123456789abcdef")
	require.NoError(t, err)

	token, _, err := maker.CreateToken("alice", "depositor", -time.Minute, AccessToken)
	require.NoError(t, err)

	_, err = maker.VerifyToken(token, AccessToken)
	require.ErrorIs(t, err, ErrExpiredToken)
}

func TestJWTInvalidKeyLength(t *testing.T) {
	_, err := NewJWTMaker("short")
	require.Error(t, err)
}

func TestJWTWrongTokenType(t *testing.T) {
	maker, _ := NewJWTMaker("0123456789abcdef0123456789abcdef")
	token, _, _ := maker.CreateToken("alice", "depositor", time.Minute, RefreshToken)

	// 拿 refresh token 去验证 access token 类型 → 应该拒绝
	_, err := maker.VerifyToken(token, AccessToken)
	require.ErrorIs(t, err, ErrInvalidToken)
}
