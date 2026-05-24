package sqlc

import (
	"context"
	"testing"

	"github.com/520wheat/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	user := createTestUser(t)
	account := createTestAccount(t, user, "USD")

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
}

func TestGetEntry(t *testing.T) {
	user := createTestUser(t)
	account := createTestAccount(t, user, "USD")

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, _ := testQueries.CreateEntry(context.Background(), arg)

	result, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.Equal(t, entry.ID, result.ID)
	require.Equal(t, entry.Amount, result.Amount)
}

func TestListEntries(t *testing.T) {
	user := createTestUser(t)
	account := createTestAccount(t, user, "USD")

	for i := 0; i < 3; i++ {
		createTestEntry(t, account)
	}

	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entries), 3)
}

func createTestEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	return entry
}
