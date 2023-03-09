package main

import (
	"fmt"
	"os"
	"strconv"

	"vlg/tools/model"

	"github.com/davecgh/go-spew/spew"
	"github.com/timshannon/badgerhold/v4"
)

const dataDirectory = "data"

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: lookup <nodeID>")
		os.Exit(1)
	}
	nodeID, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		panic(err)
	}

	options := badgerhold.DefaultOptions
	options.Dir = dataDirectory
	options.ValueDir = dataDirectory
	options.ReadOnly = true
	store, err := badgerhold.Open(options)
	defer func() {
		_ = store.Close()
	}()
	if err != nil {
		panic(err)
	}

	// retrieve the data
	record, _, err := model.RecordByID(store, nodeID)
	if err != nil {
		panic(err)
	}
	spew.Dump(record)
}
