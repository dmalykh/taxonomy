package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dmalykh/tagservice/tagservice/repository"
)

var ErrBeginxTx = errors.New(`error begin transaction`)

type TxOptions sql.TxOptions

type Transactioner interface {
	BeginTx(ctx context.Context, opts ...*TxOptions) (Transaction, error)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Tag() repository.Tag
	Namespace() repository.Namespace
	Relation() repository.Relation
	Category() repository.Category
}
