package repository

import (
	"context"
	"fmt"

	"github.com/dmalykh/tagservice/repository/entgo/ent"
	entnamespace "github.com/dmalykh/tagservice/repository/entgo/ent/namespace"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
)

func NewNamespace(client *ent.NamespaceClient) repository.Namespace {
	return &Namespace{
		client,
	}
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Namespace struct {
	client *ent.NamespaceClient
}

func (n *Namespace) Create(ctx context.Context, name string) (model.Namespace, error) {
	ns, err := n.client.Create().SetName(name).Save(ctx)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", repository.ErrCreateNamespace, err.Error())
	}

	return n.ent2model(ns), nil
}

func (n *Namespace) Update(ctx context.Context, id uint, name string) (model.Namespace, error) {
	ns, err := n.client.UpdateOneID(int(id)).SetName(name).Save(ctx)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w: %s", repository.ErrUpdateNamespace, err.Error())
	}

	return n.ent2model(ns), err
}

func (n *Namespace) GetByID(ctx context.Context, id uint) (model.Namespace, error) {
	ns, err := n.client.Get(ctx, int(id))
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w (%d): %s", repository.ErrFindNamespace, id, err.Error())
	}

	return n.ent2model(ns), err
}

func (n *Namespace) GetByName(ctx context.Context, name string) (model.Namespace, error) {
	ns, err := n.client.Query().Where(entnamespace.Name(name)).Only(ctx)
	if err != nil {
		return model.Namespace{}, fmt.Errorf("%w (%s): %s", repository.ErrFindNamespace, name, err.Error())
	}

	return n.ent2model(ns), err
}

func (n *Namespace) DeleteByID(ctx context.Context, id uint) error {
	if err := n.client.DeleteOneID(int(id)).Exec(ctx); err != nil {
		return fmt.Errorf("%w (%d): %s", repository.ErrDeleteNamespace, id, err.Error())
	}

	return nil
}

func (n *Namespace) GetList(ctx context.Context, limit, offset uint) ([]model.Namespace, error) {
	nss, err := n.client.Query().Offset(int(offset)).Limit(int(limit)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrFindNamespace, err.Error())
	}

	namespaces := make([]model.Namespace, 0, len(nss))

	for _, ns := range nss {
		namespaces = append(namespaces, n.ent2model(ns))
	}

	return namespaces, nil
}

func (n *Namespace) ent2model(ns *ent.Namespace) model.Namespace {
	return model.Namespace{
		ID:   uint(ns.ID),
		Name: ns.Name,
	}
}
