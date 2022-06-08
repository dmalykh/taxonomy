// Code generated by entc, DO NOT EDIT.

package relation

import (
	"tagservice/repository/entgo/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// TagID applies equality check predicate on the "tag_id" field. It's identical to TagIDEQ.
func TagID(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTagID), v))
	})
}

// EntityID applies equality check predicate on the "entity_id" field. It's identical to EntityIDEQ.
func EntityID(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEntityID), v))
	})
}

// NamespaceID applies equality check predicate on the "namespace_id" field. It's identical to NamespaceIDEQ.
func NamespaceID(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNamespaceID), v))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// TagIDEQ applies the EQ predicate on the "tag_id" field.
func TagIDEQ(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTagID), v))
	})
}

// TagIDNEQ applies the NEQ predicate on the "tag_id" field.
func TagIDNEQ(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTagID), v))
	})
}

// TagIDIn applies the In predicate on the "tag_id" field.
func TagIDIn(vs ...int) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTagID), v...))
	})
}

// TagIDNotIn applies the NotIn predicate on the "tag_id" field.
func TagIDNotIn(vs ...int) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTagID), v...))
	})
}

// EntityIDEQ applies the EQ predicate on the "entity_id" field.
func EntityIDEQ(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEntityID), v))
	})
}

// EntityIDNEQ applies the NEQ predicate on the "entity_id" field.
func EntityIDNEQ(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEntityID), v))
	})
}

// EntityIDIn applies the In predicate on the "entity_id" field.
func EntityIDIn(vs ...int) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldEntityID), v...))
	})
}

// EntityIDNotIn applies the NotIn predicate on the "entity_id" field.
func EntityIDNotIn(vs ...int) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldEntityID), v...))
	})
}

// EntityIDGT applies the GT predicate on the "entity_id" field.
func EntityIDGT(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldEntityID), v))
	})
}

// EntityIDGTE applies the GTE predicate on the "entity_id" field.
func EntityIDGTE(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldEntityID), v))
	})
}

// EntityIDLT applies the LT predicate on the "entity_id" field.
func EntityIDLT(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldEntityID), v))
	})
}

// EntityIDLTE applies the LTE predicate on the "entity_id" field.
func EntityIDLTE(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldEntityID), v))
	})
}

// NamespaceIDEQ applies the EQ predicate on the "namespace_id" field.
func NamespaceIDEQ(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNamespaceID), v))
	})
}

// NamespaceIDNEQ applies the NEQ predicate on the "namespace_id" field.
func NamespaceIDNEQ(v int) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldNamespaceID), v))
	})
}

// NamespaceIDIn applies the In predicate on the "namespace_id" field.
func NamespaceIDIn(vs ...int) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldNamespaceID), v...))
	})
}

// NamespaceIDNotIn applies the NotIn predicate on the "namespace_id" field.
func NamespaceIDNotIn(vs ...int) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldNamespaceID), v...))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Relation {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Relation(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// HasTag applies the HasEdge predicate on the "tag" edge.
func HasTag() predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TagTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TagTable, TagColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTagWith applies the HasEdge predicate on the "tag" edge with a given conditions (other predicates).
func HasTagWith(preds ...predicate.Tag) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TagInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TagTable, TagColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasNamespace applies the HasEdge predicate on the "namespace" edge.
func HasNamespace() predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(NamespaceTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, NamespaceTable, NamespaceColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasNamespaceWith applies the HasEdge predicate on the "namespace" edge with a given conditions (other predicates).
func HasNamespaceWith(preds ...predicate.Namespace) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(NamespaceInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, NamespaceTable, NamespaceColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Relation) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Relation) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Relation) predicate.Relation {
	return predicate.Relation(func(s *sql.Selector) {
		p(s.Not())
	})
}
