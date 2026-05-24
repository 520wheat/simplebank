package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	user := createTestUser(t)

	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: "test-refresh-token",
		UserAgent:    "Mozilla/5.0",
		ClientIp:     "127.0.0.1",
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.Username, session.Username)
}

func TestGetSession(t *testing.T) {
	user := createTestUser(t)

	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: "test-refresh-token",
		UserAgent:    "Mozilla/5.0",
		ClientIp:     "127.0.0.1",
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	session, _ := testQueries.CreateSession(context.Background(), arg)

	result, err := testQueries.GetSession(context.Background(), session.ID)
	require.NoError(t, err)
	require.Equal(t, session.ID, result.ID)
}
