package model

import (
	"fmt"
	"strings"

	"github.com/timshannon/badgerhold/v4"
)

type Relationship struct {
	FromID           string   `csv:"node_id_start"`
	FromType         string   `csv:"-"`
	ToID             string   `csv:"node_id_end"`
	ToType           string   `csv:"-"`
	RelationshipType string   `csv:"rel_type"`
	Status           string   `csv:"status"`
	StartDate        DateTime `csv:"start_date"`
	EndDate          DateTime `csv:"end_date"`
	SourceID         string   `csv:"source_id"`
}

func (r *Relationship) String() string {
	return fmt.Sprintf("%s->%s->%s", r.FromID, r.RelationshipType, r.ToID)
}

// CountRelationships counts the number of relationships of each type.
func CountRelationships(store *badgerhold.Store) string {
	var result strings.Builder
	var missingCount int
	q := &badgerhold.Query{}
	relMap := make(map[string]int)
	err := store.ForEach(q, func(record *Relationship) error {
		if record.FromType == "" || record.ToType == "" {
			missingCount++
			return nil
		}
		key := record.FromType + "->" + record.ToType + "->" + record.RelationshipType
		_, ok := relMap[key]
		if !ok {
			relMap[key] = 1
		} else {
			relMap[key]++
		}
		return nil
	})
	if err != nil {
		result.WriteString(err.Error())
	}
	result.WriteString(fmt.Sprintf("Missing: %d\n", missingCount))
	for k, v := range relMap {
		result.WriteString(fmt.Sprintf("%s: %d\n", k, v))
	}
	return result.String()
}
