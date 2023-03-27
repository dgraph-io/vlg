package model

import (
	"fmt"
	"io"
	"log"
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

/*
| connected_to             |
| intermediary_of          |
| officer_of               |
| probably_same_officer_as |
| registered_address       |
| same_address_as          |
| same_as                  |
| same_company_as          |
| same_id_as               |
| same_intermediary_as     |
| same_name_as             |
| similar                  |
| similar_company_as       |
| underlying               |
*/
func (r *Relationship) ToRDF(w io.Writer) {
	switch r.RelationshipType {
	case "connected_to":
		fmt.Fprintf(w, "_:%s <Record.connectedTo> _:%s .\n", r.FromID, r.ToID)
		fmt.Fprintf(w, "_:%s <Record.connectedTo> _:%s .\n", r.ToID, r.FromID)
	case "intermediary_of":
		fmt.Fprintf(w, "_:%s <Record.hasIntermediary> _:%s .\n", r.ToID, r.FromID)
		fmt.Fprintf(w, "_:%s <Record.intermediaryFor> _:%s .\n", r.FromID, r.ToID)
	case "officer_of":
		fmt.Fprintf(w, "_:%s <Record.hasOfficer> _:%s .\n", r.ToID, r.FromID)
		fmt.Fprintf(w, "_:%s <Record.officerFor> _:%s .\n", r.FromID, r.ToID)
	case "registered_address":
		fmt.Fprintf(w, "_:%s <Record.hasAddress> _:%s .\n", r.FromID, r.ToID)
		fmt.Fprintf(w, "_:%s <Record.addressFor> _:%s .\n", r.ToID, r.FromID)
	case "same_address_as", "same_as", "same_company_as", "same_id_as",
		"same_intermediary_as":
		fmt.Fprintf(w, "_:%s <Record.sameAs> _:%s .\n", r.FromID, r.ToID)
		fmt.Fprintf(w, "_:%s <Record.sameAs> _:%s .\n", r.ToID, r.FromID)
	case "same_name_as":
		fmt.Fprintf(w, "_:%s <Record.sameNameAs> _:%s .\n", r.FromID, r.ToID)
		fmt.Fprintf(w, "_:%s <Record.sameNameAs> _:%s .\n", r.ToID, r.FromID)
	case "probably_same_officer_as", "similar", "similar_company_as":
		fmt.Fprintf(w, "_:%s <Record.similarTo> _:%s .\n", r.FromID, r.ToID)
		fmt.Fprintf(w, "_:%s <Record.similarTo> _:%s .\n", r.ToID, r.FromID)
	default:
		log.Println("Warning, skipping relationship type", r.RelationshipType)
	}

	// TODO: write facets
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
