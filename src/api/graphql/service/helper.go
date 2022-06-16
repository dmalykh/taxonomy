package service

import (
	"github.com/AlekSi/pointer"
	"github.com/dmalykh/tagservice/api/graphql/generated/genmodel"
	apimodel "github.com/dmalykh/tagservice/api/graphql/model"
	"github.com/dmalykh/tagservice/api/graphql/service/cursor"
	"github.com/dmalykh/tagservice/tagservice/model"
	"unsafe"
)

func tag2gen(tag model.Tag) apimodel.Tag {
	return apimodel.Tag{
		ID:          int64(tag.Id),
		Name:        tag.Data.Name,
		Title:       &tag.Data.Title,
		Description: &tag.Data.Description,
		CategoryId:  int64(tag.Data.CategoryId),
	}
}
func category2gen(category model.Category) apimodel.Category {
	return apimodel.Category{
		ID:          int64(category.Id),
		Name:        category.Data.Name,
		Title:       category.Data.Title,
		Description: category.Data.Description,
		ParentId:    (*int64)(unsafe.Pointer(category.Data.ParentId)),
	}
}

func int64stoUints(ints []int64) []uint {
	var uints = make([]uint, 0, len(ints))
	for _, i := range ints {
		uints = append(uints, uint(i))
	}
	return uints
}

func tagsConnection(tags []model.Tag, limit int) *genmodel.TagsConnection {
	var connection genmodel.TagsConnection
	if len(tags) == 0 {
		return nil
	}
	if len(tags) > limit {
		tags = tags[:limit]
	}
	connection.Edges = make([]genmodel.TagsEdge, len(tags))
	for i, tag := range tags {
		var gen = tag2gen(tag)
		connection.Edges[i] = genmodel.TagsEdge{
			Cursor: cursor.Marshal(tag.Id),
			Node:   &gen,
		}
	}
	connection.PageInfo = genmodel.PageInfo{
		StartCursor: connection.Edges[0].Cursor,
		EndCursor:   connection.Edges[len(connection.Edges)-1].Cursor,
		HasNextPage: pointer.ToBool(len(tags) > limit),
	}
	return &connection
}
