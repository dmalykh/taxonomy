package taxonomy

import (
	"context"
	"github.com/dmalykh/taxonomy/taxonomy/model"
)

type Reference interface {
	// Create creates reference between specified term, namespace and all entities. The method returns an error if
	// any of reference wasn't created.
	// If reference already exists, it will be rewriting.
	Create(ctx context.Context, termID uint64, namespace string, entitiesID ...model.EntityID) error

	// Delete removes the relation between term, namespace and entities
	Delete(ctx context.Context, termID uint64, namespace string, entitiesID ...model.EntityID) error

	// Get returns all entities with specified namespace and which have all specified terms in termGroups.
	// For example:
	// 		Created *vocabularies* for laptops "RAM", "Matrix type", "Display size".
	// 		We would receive all laptops that have:
	//			- "RAM" 512 or 1024 (i.e. term's id for 512 is 92, for 1024 is 23)
	//			- "Matrix type" OLED or IPS  (i.e. term's id for OLED is 43, for IPS is 58)
	//			- "Display size" between 13 and 15 (i.e. term's id for 13 is 83, for 14 is 99, for 15 is 146)
	// 		So TermID in model.EntityFilter should given as
	//									["512", "1024"], ["OLED", "IPS"], [all terms between 13 and 26 values]:
	//			model.EntityFilter{
	//				TermID: [][]uint{
	//					{92, 23},
	//					{43, 58},
	//					{83, 99, 146},
	//				}
	//			}
	Get(ctx context.Context, filter *model.ReferenceFilter) ([]*model.Reference, error)
}
