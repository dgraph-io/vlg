package main

import (
	"log"
	"math/rand"
	"os"

	"github.com/timshannon/badgerhold/v4"

	"vlg/tools/model"
)

const exportDir = "../rdf-subset"
const exportCount = 100000

// Export a subset of the data to RDF files.
func main() {

	store, err := model.GetStore(true)
	defer func() {
		_ = store.Close()
	}()
	if err != nil {
		panic(err)
	}

	// load all relationships
	var relationships []*model.Relationship
	err = store.Find(&relationships, &badgerhold.Query{})
	if err != nil {
		panic(err)
	}
	count := len(relationships)
	f, err := os.Create(exportDir + "/data.rdf")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	// randomly select relationships (note: seed hard-coded for reproducibility)
	rand.Seed(42)
	for n := 0; n < exportCount; n++ {
		i := rand.Intn(int(count))
		from, _, err := model.RecordByID(store, relationships[i].FromID)
		if err != nil {
			panic(err)
		}
		to, _, err := model.RecordByID(store, relationships[i].ToID)
		if err != nil {
			panic(err)
		}
		from.ToRDF(f)
		to.ToRDF(f)
		relationships[i].ToRDF(f)
		if n%1000 == 0 {
			log.Printf("Exported %d relationships of %d", n, exportCount)
		}
	}
}
