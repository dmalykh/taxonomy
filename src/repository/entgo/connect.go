package entgo

import (
	"context"
	"fmt"
	"github.com/xo/dburl"
	"log"
	"tagservice/repository/entgo/ent"
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
