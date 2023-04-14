package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/timshannon/badgerhold/v4"

	"vlg/tools/model"
)

// search for records.
// example: go run search/main.go relationship 'entity->entity->officer_of'
func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: search <type> <query>")
		os.Exit(1)
	}

	store, err := model.GetStore(true)
	defer func() {
		_ = store.Close()
	}()
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "relationship":
		relationships, err := findRelationships(store, os.Args[2])
		if err != nil {
			panic(err)
		}
		for _, relationship := range relationships {
			from, _, err := model.RecordByID(store, relationship.FromID)
			if err != nil {
				panic(err)
			}
			to, _, err := model.RecordByID(store, relationship.ToID)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s %s %s\n", from, relationship.RelationshipType, to)
		}
	default:
		fmt.Println("Unknown search type: ", os.Args[1])
		os.Exit(1)
	}
}

func findRelationships(store *badgerhold.Store, query string) ([]*model.Relationship, error) {
	args := strings.Split(query, "->")
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid query: %s", query)
	}
	var relationships []*model.Relationship
	err := store.Find(&relationships, badgerhold.Where("FromType").Eq(args[0]).And("ToType").Eq(args[1]).And("RelationshipType").Eq(args[2]))
	if err != nil {
		return nil, err
	}
	return relationships, nil
}
