package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/namespace"
	"github.com/dmalykh/taxonomy/internal/repository/entgo/ent/predicate"
	"github.com/dmalykh/taxonomy/taxonomy/model"
	"github.com/dmalykh/taxonomy/taxonomy/repository"
)

func NewNamespace(client *ent.NamespaceClient) repository.Namespace {
	return &Namespace{
		client,
	}
}

type Namespace struct {
	client *ent.NamespaceClient
}

func (n *Namespace) Create(ctx context.Context, data *model.NamespaceData) (*model.Namespace, error) {
	ns, err := n.client.Create().
		SetName(data.Name).
		SetTitle(data.Title).
		Save(ctx)
	if err != nil {
		return nil, errors.Join(repository.ErrCreateNamespace, err)
	}

	return n.ent2model(ns), nil
}

func (n *Namespace) Update(ctx context.Context, id uint64, data *model.NamespaceData) (*model.Namespace, error) {
	ns, err := n.client.UpdateOneID(id).
		SetName(data.Name).
		SetTitle(data.Title).
		Save(ctx)
	if err != nil {
		return nil, errors.Join(repository.ErrUpdateNamespace, err)
	}

	return n.ent2model(ns), err
}

func (n *Namespace) Delete(ctx context.Context, filter *repository.NamespaceFilter) error {
	_, err := n.client.Delete().Where(
		n.buildQuery(filter)...,
	).Exec(ctx)
	if err != nil {
		return errors.Join(repository.ErrDeleteNamespace, err)
	}

	return nil
}

func (n *Namespace) Get(ctx context.Context, filter *repository.NamespaceFilter) ([]*model.Namespace, error) {
	nss, err := n.client.Query().Where(
		n.buildQuery(filter)...,
	).Limit(int(filter.Limit)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrFindNamespace, err.Error())
	}

	namespaces := make([]*model.Namespace, 0, len(nss))

	for _, ns := range nss {
		namespaces = append(namespaces, n.ent2model(ns))
	}

	return namespaces, nil
}

func (n *Namespace) buildQuery(filter *repository.NamespaceFilter) []predicate.Namespace {
	var predicates = make([]predicate.Namespace, 0)

	// Filter by fill name
	if len(filter.Name) > 0 {
		predicates = append(predicates, namespace.NameIn(filter.Name...))
	}

	// After id condition
	if filter.AfterID != nil {
		predicates = append(predicates, namespace.IDGT(*filter.AfterID))
	}

	return predicates
}

func (n *Namespace) ent2model(ns *ent.Namespace) *model.Namespace {
	return &model.Namespace{
		ID: ns.ID,
		Data: model.NamespaceData{
			Name:  ns.Name,
			Title: ns.Title,
		},
	}
}
