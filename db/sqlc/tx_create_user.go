package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserResult is the result of the transfer transaction
type CreateUserResult struct {
	User User
}

func (store *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserResult, error) {
	var result CreateUserResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})

	return result, err
}
