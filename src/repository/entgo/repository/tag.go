package repository

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"tagservice/repository/entgo/ent"
	"tagservice/repository/entgo/ent/tag"
	"tagservice/server/model"
	"tagservice/server/repository"
)

var (
	ErrCreateTag = errors.New(`failed to create tag`)
	ErrUpdateTag = errors.New(`failed to update tag`)
	ErrFindTag   = errors.New(`failed to find tag`)
	ErrDeleteTag = errors.New(`failed to delete tag`)
)

type Tag struct {
	client *ent.TagClient
}

func NewTag(client *ent.TagClient) repository.Tag {
	return &Tag{
		client: client,
	}
}

func (t *Tag) Create(ctx context.Context, data *model.TagData) (model.Tag, error) {
	ns, err := t.client.Create().
		SetName(data.Name).
		SetTitle(data.Title).
		SetDescription(data.Description).
		SetCategoryID(int(data.CategoryId)).
		Save(ctx)

	if err != nil {
		return model.Tag{}, fmt.Errorf("%w: %s", ErrCreateTag, err.Error())
	}
	return t.ent2model(ns), nil
}

func (t *Tag) Update(ctx context.Context, id uint, data *model.TagData) (model.Tag, error) {
	tag, err := t.client.UpdateOneID(int(id)).
		SetName(data.Name).
		SetTitle(data.Title).
		SetDescription(data.Description).
		SetCategoryID(int(data.CategoryId)).
		Save(ctx)
	if err != nil {
		return model.Tag{}, fmt.Errorf("%w: %s", ErrUpdateTag, err.Error())
	}
	return t.ent2model(tag), err
}

func (t *Tag) DeleteById(ctx context.Context, id uint) error {
	if err := t.client.DeleteOneID(int(id)).Exec(ctx); err != nil {
		return fmt.Errorf("%w (%d): %s", ErrDeleteTag, id, err.Error())
	}
	return nil
}

func (t *Tag) GetById(ctx context.Context, id uint) (model.Tag, error) {
	ns, err := t.client.Get(ctx, int(id))
	if err != nil {
		return model.Tag{}, fmt.Errorf("%w (%d): %s", ErrFindTag, id, err.Error())
	}
	return t.ent2model(ns), err
}

func (t *Tag) GetList(ctx context.Context, limit, offset uint) ([]model.Tag, error) {
	enttags, err := t.client.Query().Offset(int(offset)).Limit(int(limit)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFindTag, err.Error())
	}
	var tags = make([]model.Tag, 0, len(enttags))
	for _, enttag := range enttags {
		tags = append(tags, t.ent2model(enttag))
	}
	return tags, nil
}

func (t *Tag) ent2model(tag *ent.Tag) model.Tag {
	return model.Tag{
		Id: uint(tag.ID),
		Data: model.TagData{
			Name:        tag.Name,
			Title:       tag.Title,
			Description: tag.Description,
			CategoryId:  uint(tag.CategoryID),
		},
	}
}

func (t *Tag) GetByName(ctx context.Context, name string) (model.Tag, error) {
	enttag, err := t.client.Query().Where(tag.Name(name)).Only(ctx)
	if err != nil {
		return model.Tag{}, fmt.Errorf("%w (%s): %s", ErrFindTag, name, err.Error())
	}
	return t.ent2model(enttag), err
}

func (t *Tag) GetByFilter(ctx context.Context, filter model.TagFilter, limit, offset uint) ([]model.Tag, error) {
	enttags, err := t.client.Query().Where(func(s *sql.Selector) {
		// Filter by categories id
		if len(filter.CategoryId) > 0 {
			s.Where(sql.InInts(tag.CategoryColumn, func() []int {
				var ints = make([]int, len(filter.CategoryId))
				for i, val := range filter.CategoryId {
					ints[i] = int(val)
				}
				return ints
			}()...))
		}
	}).Limit(int(limit)).Offset(int(offset)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFindTag, err.Error())
	}
	var tags = make([]model.Tag, 0, len(enttags))
	for _, enttag := range enttags {
		tags = append(tags, t.ent2model(enttag))
	}
	return tags, nil
}
