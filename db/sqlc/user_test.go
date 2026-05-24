package sqlc

import (
	"context"
	"testing"

	"github.com/520wheat/simplebank/util"
	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateUser(t *testing.T) {
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
		Role:           "depositor",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.False(t, user.IsEmailVerified)
}

func TestGetUser(t *testing.T) {
	user := createTestUser(t)

	result, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, result.Username)
	require.Equal(t, user.Email, result.Email)
}

func TestUpdateUser(t *testing.T) {
	user := createTestUser(t)

	newFullName := util.RandomFullName()
	result, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		FullName: pgtype.Text{String: newFullName, Valid: true},
	})
	require.NoError(t, err)
	require.Equal(t, newFullName, result.FullName)
	require.Equal(t, user.HashedPassword, result.HashedPassword) // 没传，不变
}
