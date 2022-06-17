package transaction

import (
	"context"
	"fmt"
	"log"

	"github.com/dmalykh/tagservice/repository/entgo/ent"
	entrepo "github.com/dmalykh/tagservice/repository/entgo/repository"
	"github.com/dmalykh/tagservice/tagservice/repository"
	"github.com/dmalykh/tagservice/tagservice/repository/transaction"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Transaction struct {
	tx *ent.Tx
	ns repository.Namespace
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Transactioner struct {
	client *ent.Client
}

func New(client *ent.Client) transaction.Transactioner {
	return &Transactioner{
		client: client,
	}
}

func (t *Transactioner) BeginTx(ctx context.Context, _ ...*transaction.TxOptions) (transaction.Transaction, error) {
	tx, err := t.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf(`%w: %s`, transaction.ErrBeginxTx, err.Error())
	}
	// Rollback transaction immediately when context done
	go func() {
		<-ctx.Done()

		err := tx.Rollback()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	newTx := Transaction{
		tx: tx,
	}

	return &newTx, nil
}

func (t *Transaction) Commit(_ context.Context) error {
	return t.tx.Commit() //nolint:wrapcheck
}

func (t *Transaction) Rollback(_ context.Context) error {
	return t.tx.Rollback() //nolint:wrapcheck
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
