package tagservice

import (
	"context"
	"errors"

	"github.com/dmalykh/tagservice/tagservice/model"
)

var (
	ErrNamespaceNotFound   = errors.New(`namespace not found`)
	ErrNamespaceNotCreated = errors.New(`namespace have not created`)
	ErrNamespaceNotUpdated = errors.New(`namespace have not updated`)
)

type Namespace interface {
	Create(ctx context.Context, name string) (model.Namespace, error)
	Update(ctx context.Context, uid uint, name string) (model.Namespace, error)
	Delete(ctx context.Context, uid uint) error
	GetList(ctx context.Context, limit, offset uint) ([]model.Namespace, error)
	GetByName(ctx context.Context, namespace string) (model.Namespace, error)
}
