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
	f, err := os.Create(exportDir + "/data.rdf")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	count := len(relationships)
	relationshipSet := make(map[int]struct{})
	// randomly select relationships (note: seed hard-coded for reproducibility)
	rand.Seed(42)
	for n := 0; n < exportCount; n++ {
		i := rand.Intn(int(count))
		if _, ok := relationshipSet[i]; ok {
			n-- // try again
			continue
		}
		relationshipSet[i] = struct{}{}
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
		if n%10000 == 0 {
			log.Printf("Exporting %d of %d relationships ...", n, exportCount)
		}
	}
	log.Println("Exported", len(relationshipSet), "random relationships")
	// load some notable records
	notables := map[string]struct{}{
		"182471":   {}, // Imee Marcos: The daughter of former Philippine dictator Ferdinand Marcos
		"80133535": {}, // Rex Tillerson: Former ExxonMobil CEO and U.S. Secretary of State
		"80118878": {}, // Wilbur Ross: Former U.S. Secretary of Commerce
		"84100003": {}, // The Duchy of Lancaster: The private estate of the British monarch
		"56073968": {}, // Shakira: Colombian singer
	}
	notableCount := 0
	for i := range relationships {
		r := relationships[i]
		_, fromFound := notables[r.FromID]
		_, toFound := notables[r.ToID]
		if toFound || fromFound {
			from, _, err := model.RecordByID(store, r.FromID)
			if err != nil {
				panic(err)
			}
			to, _, err := model.RecordByID(store, r.ToID)
			if err != nil {
				panic(err)
			}
			from.ToRDF(f)
			to.ToRDF(f)
			r.ToRDF(f)
			notableCount++
		}
	}
	log.Printf("Exported %d notable relationships", notableCount)
}
