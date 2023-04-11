package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"vlg/tools/model"
)

// Lookup retrieves a record from the badger store by ID.
func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: lookup <nodeID>")
		os.Exit(1)
	}
	_, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic(err)
	}

	store, err := model.GetStore(true)
	defer func() {
		_ = store.Close()
	}()
	if err != nil {
		panic(err)
	}

	// retrieve the data
	record, recordType, err := model.RecordByID(store, os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println("Found record of type", recordType)
	b, _ := json.MarshalIndent(record, "", "  ")
	fmt.Println(string(b))
}
