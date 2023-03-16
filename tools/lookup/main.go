package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"

	"vlg/tools/model"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: lookup <nodeID>")
		os.Exit(1)
	}
	nodeID, err := strconv.ParseInt(os.Args[1], 10, 64)
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
	record, _, err := model.RecordByID(store, nodeID)
	if err != nil {
		panic(err)
	}
	spew.Dump(record)
}
