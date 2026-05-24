package sqlc

import (
	"context"
	"testing"

	"github.com/520wheat/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	// 1. 先创建用户——因为 account.owner 引用 users.username
	user := createTestUser(t)

	// 2. 创建账户
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, account.ID)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
}

func TestGetAccount(t *testing.T) {
	user := createTestUser(t)
	account := createTestAccount(t, user, "USD")

	result, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.ID, result.ID)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, account.Balance, result.Balance)
	require.Equal(t, account.Currency, result.Currency)
}

func TestListAccounts(t *testing.T) {
	user := createTestUser(t)
	currencies := []string{"USD", "EUR", "CNY", "JPY", "GBP"}
	for _, c := range currencies {
		createTestAccount(t, user, c)
	}

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Owner:  user.Username,
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)
}

func TestAddAccountBalance(t *testing.T) {
	user := createTestUser(t)
	account := createTestAccount(t, user, "USD")

	arg := AddAccountBalanceParams{
		ID:     account.ID,
		Amount: 100,
	}

	result, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account.Balance+100, result.Balance)
}

func TestDeleteAccount(t *testing.T) {
	user := createTestUser(t)
	account := createTestAccount(t, user, "USD")

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	_, err = testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err) // 确认真的删掉了
}

func createTestUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "test-hash",
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
		Role:           "depositor",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, user.ID)
	return user
}

func createTestAccount(t *testing.T, user User, currency string) Account {
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: currency,
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, account.ID)
	return account
}
