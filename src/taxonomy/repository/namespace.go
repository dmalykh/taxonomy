package repository

import (
	"context"
	"errors"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

var (
	ErrCreateNamespace = errors.New(`failed to create namespace`)
	ErrUpdateNamespace = errors.New(`failed to update namespace`)
	ErrFindNamespace   = errors.New(`failed to find namespace`)
	ErrDeleteNamespace = errors.New(`failed to delete namespace`)
)

type Namespace interface {
	Create(ctx context.Context, data *model.NamespaceData) (*model.Namespace, error)
	Update(ctx context.Context, id uint64, data *model.NamespaceData) (*model.Namespace, error)
	Delete(ctx context.Context, filter *NamespaceFilter) error
	Get(ctx context.Context, filter *NamespaceFilter) ([]*model.Namespace, error)
}

type NamespaceFilter struct {
	ID      []uint64
	Name    []string
	AfterID *uint64
	Limit   uint
}
