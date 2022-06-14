package tagservice

import (
	"context"
	"errors"
	"github.com/dmalykh/tagservice/tagservice/model"
)

var (
	ErrTagNotFound           = errors.New(`tag not found`)
	ErrTagNamespaceNotFound  = errors.New(`tag's namespace not found`)
	ErrTagNotCreated         = errors.New(`tag had not created`)
	ErrTagRelationNotCreated = errors.New(`tag's relation had not created`)
	ErrTagRelationNotRemoved = errors.New(`tag's relation had not removed`)
	ErrTagNotUpdated         = errors.New(`tag have not updated`)
)

type Tag interface {
	Create(ctx context.Context, data *model.TagData) (model.Tag, error)
	Update(ctx context.Context, id uint, data *model.TagData) (model.Tag, error)
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (model.Tag, error)
	GetByName(ctx context.Context, name string, categoryId uint) (model.Tag, error)

	// GetList returns slice with tags that proper for conditions. Set nil category_id to receive tags from all categories.
	GetList(ctx context.Context, categoryId uint, limit, offset uint) ([]model.Tag, error)

	// SetRelation create relation between specified tag, namespace and all entities. Return error if any of relation didn't create.
	// If relation already exists, it will be rewriting.
	SetRelation(ctx context.Context, tagId uint, namespace string, entitiesId ...uint) error
	UnsetRelation(ctx context.Context, tagId uint, namespace string, entitiesId ...uint) error

	// GetRelationEntities returns all entities with specified namespace and which have all specified tags in tagGroups. All tags
	// specified in one tagGroup use "OR" operand, between tagGroups "AND" operand used.
	// For example:
	// 		Created categories for laptops "RAM", "Matrix type", "Display size".
	// 		We would receive all laptops that have: "RAM" (512 or 1024) and "Matrix type" (OLED or IPS) and "Display size" (between 13 and 16)
	// Use tagGroups, you should previously receive id of desirable tags, for example used names:
	//		["512", "1024"], ["OLED", "IPS"], [all tags between 13 and 26 values]
	GetRelationEntities(ctx context.Context, namespaceName string, tagGroups [][]uint) ([]model.Relation, error)
	GetTagsByEntities(ctx context.Context, namespaceName string, entities ...uint) ([]model.Tag, error)
}
