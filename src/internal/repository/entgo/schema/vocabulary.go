package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Vocabulary holds the schema definition for the Vocabulary entity.
type Vocabulary struct {
	ent.Schema
}

// Fields of the Vocabulary.
func (Vocabulary) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64(`id`).Immutable(),
		field.String(`name`).NotEmpty(),
		field.String(`title`).Optional(),
		field.Text(`description`).Optional(),
		field.Uint64(`parent_id`).Optional().Nillable(),
	}
}

func (Vocabulary) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`name`),
		index.Fields(`parent_id`),
		index.Fields(`name`, `parent_id`).Unique(),
	}
}

// Edges of the Vocabulary.
func (Vocabulary) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(`term`, Term.Type),
		edge.To(`children`, Vocabulary.Type),
		edge.From(`parent`, Vocabulary.Type).
			Ref(`children`),
	}
}
