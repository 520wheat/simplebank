package sqlc

import (
	"context"
	"testing"

	"github.com/520wheat/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	user := createTestUser(t)
	fromAccount := createTestAccount(t, user, "USD")
	toAccount := createTestAccount(t, user, "EUR")

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, transfer.ID)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
}

func TestGetTransfer(t *testing.T) {
	user := createTestUser(t)
	from := createTestAccount(t, user, "USD")
	to := createTestAccount(t, user, "EUR")
	transfer := createTestTransfer(t, from, to)

	result, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, result.ID)
}

func TestListTransfers(t *testing.T) {
	user := createTestUser(t)
	from := createTestAccount(t, user, "USD")
	to := createTestAccount(t, user, "EUR")

	for i := 0; i < 3; i++ {
		createTestTransfer(t, from, to)
	}

	transfers, err := testQueries.ListTransfers(context.Background(),
		ListTransfersParams{
			FromAccountID: from.ID,
			ToAccountID:   from.ID,
			Limit:         5,
			Offset:        0,
		})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(transfers), 3)
}

func createTestTransfer(t *testing.T, from, to Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	return transfer
}
