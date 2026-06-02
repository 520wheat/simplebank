package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
)

// VerifyEmailTxParams 邮箱验证事务的参数
type VerifyEmailTxParams struct {
	EmailId    int64  `json:"email_id"`
	SecretCode string `json:"secret_code"`
}

// VerifyEmailTxResult 邮箱验证事务的返回
type VerifyEmailTxResult struct {
	User        User        `json:"user"`
	VerifyEmail VerifyEmail `json:"verify_email"`
}

// VerifyEmailTx 在事务中验证邮箱：标记验证码已使用 + 标记用户邮箱已验证
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. 标记验证码已使用
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		// 2. 标记用户邮箱已验证
		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}