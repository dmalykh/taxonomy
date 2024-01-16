package taxonomy

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

var (
	ErrNamespaceNotFound   = errors.New(`namespace not found`)
	ErrNamespaceNotCreated = errors.New(`namespace have not created`)
	ErrNamespaceNotUpdated = errors.New(`namespace have not updated`)
	ErrNamespaceNotDeleted = errors.New(`namespace have not deleted`)
)

type Namespace interface {
	Create(ctx context.Context, name string) (*model.Namespace, error)
	Update(ctx context.Context, uid uint64, name string) (*model.Namespace, error)
	Delete(ctx context.Context, uid uint64) error
	GetByName(ctx context.Context, namespace string) (*model.Namespace, error)

	Get(ctx context.Context, limit uint, afterId *uint64) ([]*model.Namespace, error)
}
