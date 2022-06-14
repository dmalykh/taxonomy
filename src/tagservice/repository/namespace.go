package repository

import (
	"context"
	"errors"
	"github.com/dmalykh/tagservice/tagservice/model"
)

var (
	ErrCreateNamespace = errors.New(`failed to create namespace`)
	ErrUpdateNamespace = errors.New(`failed to update namespace`)
	ErrFindNamespace   = errors.New(`failed to find namespace`)
	ErrDeleteNamespace = errors.New(`failed to delete namespace`)
)

type Namespace interface {
	Create(ctx context.Context, name string) (model.Namespace, error)
	Update(ctx context.Context, id uint, name string) (model.Namespace, error)
	GetById(ctx context.Context, id uint) (model.Namespace, error)
	GetByName(ctx context.Context, name string) (model.Namespace, error)
	DeleteById(ctx context.Context, id uint) error
	GetList(ctx context.Context, limit, offset uint) ([]model.Namespace, error)
}
