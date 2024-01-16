package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Namespace holds the schema definition for the Namespace entity.
type Namespace struct {
	ent.Schema
}

// Fields of the Namespace.
func (Namespace) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64(`id`).Immutable(),
		field.String(`name`).NotEmpty().Unique(),
		field.String(`title`).Optional(),
	}
}

func (Namespace) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`name`).Unique(),
	}
}

// Edges of the Namespace.
func (Namespace) Edges() []ent.Edge {
	return nil
}
