package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// Relation holds the schema definition for the Relation entity.
type Relation struct {
	ent.Schema
}

// Fields of the Relation.
func (Relation) Fields() []ent.Field {
	return []ent.Field{
		field.Int(`tag_id`).Positive(),
		field.Int(`entity_id`).Positive(),
		field.Int(`namespace_id`).Positive(),
		field.Time(`created_at`).Immutable().Default(time.Now),
	}
}

func (Relation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`tag_id`, `namespace_id`, `entity_id`).Unique(),
		index.Fields(`tag_id`, `namespace_id`),
		index.Fields(`tag_id`),
		index.Fields(`entity_id`),
		index.Fields(`namespace_id`),
	}
}

// Edges of the Tag.
func (Relation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(`tag`, Tag.Type).Field(`tag_id`).Unique().Required(),
		edge.To(`namespace`, Namespace.Type).Field(`namespace_id`).Unique().Required(),
	}
}
