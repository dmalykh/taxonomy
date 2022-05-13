package transaction

import (
	"context"
	"database/sql"
)

type TxOptions sql.TxOptions

type Transaction interface {
	BeginTx(ctx context.Context, opts ...*TxOptions) (Transaction, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
