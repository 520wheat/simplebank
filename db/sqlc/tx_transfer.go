package sqlc

import "context"

// TransferTxParams 转账事务的输入参数
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult 转账事务的返回结果
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx 执行一次转账。它在一个数据库事务中完成：
// 1. 创建 transfer 记录
// 2. 创建 from_entry（出账，金额为负）
// 3. 创建 to_entry（入账，金额为正）
// 4. from 账户扣款
// 5. to 账户加款
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult,
	error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. 创建转账记录
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2. 出账分录
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 3. 入账分录
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// 4-5. 更新余额。按 ID 排序防止死锁
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = updateBalances(
				ctx, q,
				arg.FromAccountID, -arg.Amount,
				arg.ToAccountID, arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = updateBalances(
				ctx, q,
				arg.ToAccountID, arg.Amount,
				arg.FromAccountID, -arg.Amount,
			)
		}

		return err
	})

	return result, err
}

func updateBalances(
	ctx context.Context,
	q *Queries,
	accountID1 int64, amount1 int64,
	accountID2 int64, amount2 int64,
) (Account, Account, error) {
	account1, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	account2, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return account1, account2, err
}
