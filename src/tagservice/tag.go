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
	GetByID(ctx context.Context, id uint) (model.Tag, error)

	// GetList returns slice with tags that proper for conditions. Set nil category_id to receive tags from all categories.
	GetList(ctx context.Context, filter *model.TagFilter) ([]model.Tag, error)

	// SetRelation create relation between specified tag, namespace and all entities. Return error if any of relation
	// didn't create.
	// If relation already exists, it will be rewriting.
	SetRelation(ctx context.Context, tagID uint, namespace string, entitiesID ...uint) error
	UnsetRelation(ctx context.Context, tagID uint, namespace string, entitiesID ...uint) error

	// GetRelations returns all entities with specified namespace and which have all specified tags in tagGroups.
	// For example:
	// 		Created *categories* for laptops "RAM", "Matrix type", "Display size".
	// 		We would receive all laptops that have:
	//			- "RAM" 512 or 1024 (i.e. tag's id for 512 is 92, for 1024 is 23)
	//			- "Matrix type" OLED or IPS  (i.e. tag's id for OLED is 43, for IPS is 58)
	//			- "Display size" between 13 and 15 (i.e. tag's id for 13 is 83, for 14 is 99, for 15 is 146)
	// 		So TagID in model.EntityFilter should given as
	//									["512", "1024"], ["OLED", "IPS"], [all tags between 13 and 26 values]:
	//			model.EntityFilter{
	//				TagID: [][]uint{
	//					{92, 23},
	//					{43, 58},
	//					{83, 99, 146},
	//				}
	//			}
	GetRelations(ctx context.Context, filter *model.EntityFilter) ([]model.Relation, error)
	GetTagsByEntities(ctx context.Context, namespaceName string, entities ...uint) ([]model.Tag, error)
}
