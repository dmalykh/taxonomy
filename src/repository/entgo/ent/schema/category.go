package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String(`name`).NotEmpty(),
		field.String(`title`).Optional(),
		field.Text(`description`).Optional(),
		field.Int(`parent_id`).Optional().Nillable(),
	}
}

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`name`),
		index.Fields(`parent_id`),
		index.Fields(`name`, `parent_id`).Unique(),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(`children`, Category.Type),
		edge.From(`parent`, Category.Type).
			Ref(`children`),
	}
}
