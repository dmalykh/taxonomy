package repository

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/dmalykh/tagservice/repository/entgo/ent"
	"github.com/dmalykh/tagservice/repository/entgo/ent/tag"
	"github.com/dmalykh/tagservice/tagservice/model"
	"github.com/dmalykh/tagservice/tagservice/repository"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
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
		SetCategoryID(int(data.CategoryID)).
		Save(ctx)
	if err != nil {
		return model.Tag{}, fmt.Errorf("%w: %s", repository.ErrCreateTag, err.Error())
	}

	return t.ent2model(ns), nil
}

func (t *Tag) Update(ctx context.Context, id uint, data *model.TagData) (model.Tag, error) {
	updated, err := t.client.UpdateOneID(int(id)).
		SetName(data.Name).
		SetTitle(data.Title).
		SetDescription(data.Description).
		SetCategoryID(int(data.CategoryID)).
		Save(ctx)
	if err != nil {
		return model.Tag{}, fmt.Errorf("%w: %s", repository.ErrUpdateTag, err.Error())
	}

	return t.ent2model(updated), err
}

func (t *Tag) DeleteByID(ctx context.Context, id uint) error {
	if err := t.client.DeleteOneID(int(id)).Exec(ctx); err != nil {
		return fmt.Errorf("%w (%d): %s", repository.ErrDeleteTag, id, err.Error())
	}

	return nil
}

func (t *Tag) GetByID(ctx context.Context, id uint) (model.Tag, error) {
	ns, err := t.client.Get(ctx, int(id))
	if err != nil {
		return model.Tag{}, fmt.Errorf("%w (%d): %s", repository.ErrFindTag, id, err.Error())
	}

	return t.ent2model(ns), err
}

func (t *Tag) ent2model(tag *ent.Tag) model.Tag {
	return model.Tag{
		ID: uint(tag.ID),
		Data: model.TagData{
			Name:        tag.Name,
			Title:       tag.Title,
			Description: tag.Description,
			CategoryID:  uint(tag.CategoryID),
		},
	}
}

func (t *Tag) GetByName(ctx context.Context, name string) ([]model.Tag, error) {
	enttags, err := t.client.Query().Where(tag.Name(name)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w (%s): %s", repository.ErrFindTag, name, err.Error())
	}

	tags := make([]model.Tag, 0, len(enttags))

	for _, enttag := range enttags {
		tags = append(tags, t.ent2model(enttag))
	}

	return tags, nil
}

func (t *Tag) GetList(ctx context.Context, filter *model.TagFilter) ([]model.Tag, error) {
	query := t.client.Query().Where(func(s *sql.Selector) {
		// Filter by categories id
		if len(filter.CategoryID) > 0 {
			s.Where(sql.InInts(tag.FieldCategoryID, func() []int {
				ints := make([]int, len(filter.CategoryID))
				for i, val := range filter.CategoryID {
					ints[i] = int(val)
				}

				return ints
			}()...))
		}
		// Filter by name
		if filter.Name != nil {
			s.Where(sql.EQ(tag.FieldName, *filter.Name))
		}
		// Add condition to id
		if filter.AfterID != nil {
			s.Where(sql.GT(tag.FieldID, *filter.AfterID))
		}
	})

	if filter.Limit != 0 || filter.Offset != 0 {
		query = query.Limit(int(filter.Limit)).Offset(int(filter.Offset))
	}

	enttags, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrFindTag, err.Error())
	}

	tags := make([]model.Tag, 0, len(enttags))

	for _, enttag := range enttags {
		tags = append(tags, t.ent2model(enttag))
	}

	return tags, nil
}
