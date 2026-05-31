package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testConnPool)

	user := createTestUser(t)
	from := createTestAccount(t, user, "USD")
	to := createTestAccount(t, user, "EUR")

	// 从账户存入 100 防止余额不足
	_, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID: from.ID, Amount: 100,
	})
	require.NoError(t, err)
	from, err = testQueries.GetAccount(context.Background(), from.ID)
	require.NoError(t, err)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// 并发执行 n 次相同方向的转账
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: from.ID,
				ToAccountID:   to.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// 收集结果
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// 验证 transfer
		require.Equal(t, from.ID, result.Transfer.FromAccountID)
		require.Equal(t, to.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = testQueries.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// 验证 entry
		require.Equal(t, from.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)

		require.Equal(t, to.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)

		// 验证最终余额
		fromAccount := result.FromAccount
		require.Equal(t, from.ID, fromAccount.ID)
		require.NotZero(t, fromAccount.Balance)

		toAccount := result.ToAccount
		require.Equal(t, to.ID, toAccount.ID)
		require.NotZero(t, toAccount.Balance)

		// 验证：from.old - amount * (i+1) = from.new
		diff1 := fromAccount.Balance - from.Balance
		diff2 := toAccount.Balance - to.Balance
		require.Equal(t, -amount*int64(i+1), diff1)
		require.Equal(t, amount*int64(i+1), diff2)

		// 确保每次 transfer ID 不重复
		k := int(result.Transfer.ID)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testConnPool)

	user := createTestUser(t)
	account1 := createTestAccount(t, user, "USD")
	account2 := createTestAccount(t, user, "EUR")

	// 两个账户都存入余额
	addBalance(t, account1.ID, 100)
	addBalance(t, account2.ID, 100)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	// 一半的并发从 1→2，另一半 2→1，故意制造潜在死锁
	for i := 0; i < n; i++ {
		from := account1
		to := account2
		if i%2 == 1 {
			from = account2
			to = account1
		}

		go func(from, to Account) {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: from.ID,
				ToAccountID:   to.ID,
				Amount:        amount,
			})
			errs <- err
		}(from, to)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
}

func addBalance(t *testing.T, accountID int64, amount int64) {
	_, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID: accountID, Amount: amount,
	})
	require.NoError(t, err)
}
