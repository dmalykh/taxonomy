package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Reference holds the schema definition for the Reference entity.
type Reference struct {
	ent.Schema
}

// Fields of the Reference.
func (Reference) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64(`id`).Immutable(),
		field.Uint64(`term_id`).Positive(),
		field.String(`entity_id`).NotEmpty().NotEmpty(),
		field.Uint64(`namespace_id`).Positive(),
		field.Time(`created_at`).Immutable().Default(time.Now),
	}
}

func (Reference) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`term_id`, `namespace_id`, `entity_id`).Unique(),
		index.Fields(`term_id`, `namespace_id`),
		index.Fields(`term_id`),
		index.Fields(`entity_id`),
		index.Fields(`namespace_id`),
	}
}

// Edges of the Term.
func (Reference) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(`term`, Term.Type).Field(`term_id`).Required().Unique(),
		edge.To(`namespace`, Namespace.Type).Field(`namespace_id`).Unique().Required(),
	}
}
