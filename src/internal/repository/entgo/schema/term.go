package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Term holds the schema definition for the Term entity.
type Term struct {
	ent.Schema
}

// Fields of the Term.
func (Term) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64(`id`).Immutable(),
		field.String(`name`).NotEmpty(),
		field.String(`title`).Optional(),
		field.Text(`description`).Optional(),
	}
}

func (Term) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`name`),
	}
}

// Edges of the Term.
func (Term) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From(`vocabulary`, Vocabulary.Type).
			Ref(`term`).Required(),

		edge.To("subterms", Term.Type).
			From("superterms"),
	}
}
