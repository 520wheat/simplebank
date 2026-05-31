package sqlc

import (
	"context"
	"fmt"
)

// execTx 在一个数据库事务中执行传入的函数
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(context.Background())

	q := New(store.connPool)
	q = q.WithTx(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(context.Background()); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(context.Background())
}
