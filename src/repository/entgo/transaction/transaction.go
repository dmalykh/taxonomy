package transaction

import (
	"context"
	"tagservice/repository/entgo/ent"
	entrepo "tagservice/repository/entgo/repository"
	"tagservice/server/repository"
	"tagservice/server/repository/transaction"
)

type Transaction struct {
	tx *ent.Tx
	ns repository.Namespace
}

type Transactioner struct {
	client *ent.Client
}

//@TODO panic, recovery, ctx.Done
func (t *Transactioner) BeginTx(ctx context.Context, opts ...*transaction.TxOptions) (transaction.Transaction, error) {
	tx, err := t.client.Tx(ctx)
	if err != nil {
		return nil, err //@TODO
	}
	var newTx = Transaction{
		tx: tx,
	}
	return &newTx, nil

}

func (t *Transaction) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *Transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

func (t *Transaction) Tag() repository.Tag {
	return entrepo.NewTag(t.tx.Tag)
}

func (t *Transaction) Namespace() repository.Namespace {
	return entrepo.NewNamespace(t.tx.Namespace)
}

func (t *Transaction) Category() repository.Category {
	return entrepo.NewCategory(t.tx.Category)
}

func (t *Transaction) Relation() repository.Relation {
	return entrepo.NewRelation(t.tx.Relation)
}
