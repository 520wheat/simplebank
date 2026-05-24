package sqlc

import (
	"context"
	"testing"

	"github.com/520wheat/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateVerifyEmail(t *testing.T) {
	user := createTestUser(t)

	arg := CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}

	ve, err := testQueries.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, ve.ID)
	require.Equal(t, arg.Username, ve.Username)
	require.False(t, ve.IsUsed)
}

func TestGetVerifyEmail(t *testing.T) {
	user := createTestUser(t)
	ve := createTestVerifyEmail(t, user)

	result, err := testQueries.GetVerifyEmail(context.Background(), ve.ID)
	require.NoError(t, err)
	require.Equal(t, ve.ID, result.ID)
}

func TestUpdateVerifyEmail(t *testing.T) {
	user := createTestUser(t)
	ve := createTestVerifyEmail(t, user)

	result, err := testQueries.UpdateVerifyEmail(context.Background(),
		UpdateVerifyEmailParams{
			ID:         ve.ID,
			SecretCode: ve.SecretCode,
		})
	require.NoError(t, err)
	require.True(t, result.IsUsed)
}

func createTestVerifyEmail(t *testing.T, user User) VerifyEmail {
	arg := CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}
	ve, err := testQueries.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	return ve
}
