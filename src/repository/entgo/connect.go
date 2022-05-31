package entgo

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xo/dburl"
	"log"
	"tagservice/repository/entgo/ent"
	tx "tagservice/repository/entgo/transaction"
	"tagservice/server/repository/transaction"
)

func Connect(ctx context.Context, dsn string) (*ent.Client, error) {
	u, err := dburl.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf(`wrong dsn %w`, err)
	}

	client, err := ent.Open(u.Driver, u.DSN)
	if err != nil {
		return nil, fmt.Errorf(`(sql open "%s") %w`, u.Driver, err)
	}

	//Shutdown database connection
	go func() {
		<-ctx.Done()
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	return client, nil
}

func Init(ctx context.Context, dsn string) error {
	client, err := Connect(ctx, dsn)
	if err != nil {
		return err
	}

	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		return err
	}
	return nil
}

func Transactioner(client *ent.Client) transaction.Transactioner {
	return tx.New(client)
}
