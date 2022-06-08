package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String(`name`).NotEmpty(),
		field.String(`title`).Optional(),
		field.Text(`description`).Optional(),
		field.Int(`category_id`).NonNegative(),
	}
}

func (Tag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`category_id`),
		index.Fields(`name`, `category_id`).Unique(),
	}
}

// Edges of the Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(`category`, Category.Type).
			Unique().Required().
			Field("category_id"),
	}
}
