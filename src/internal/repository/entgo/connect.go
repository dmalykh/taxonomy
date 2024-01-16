//nolint:revive
package entgo

import (
	"context"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/migrate"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/xiaoqidun/entps"
	"github.com/xo/dburl"
)

func Connect(ctx context.Context, dsn string, debug bool) (*ent.Client, error) {
	u, err := dburl.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf(`wrong dsn %w`, err)
	}

	client, err := ent.Open(u.Driver, u.DSN)
	if err != nil {
		return nil, fmt.Errorf(`(sql open "%s") %w`, u.Driver, err)
	}

	// Shutdown database connection.
	go func() {
		<-ctx.Done()

		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if debug {
		client = client.Debug()
	}

	return client, nil
}

func Init(ctx context.Context, dsn string, verbose bool) error {
	client, err := Connect(ctx, dsn, verbose)
	if err != nil {
		return err
	}

	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true), schema.WithAtlas(true)); err != nil {
		return fmt.Errorf(`error create schema: %w`, err)
	}

	return nil
}
