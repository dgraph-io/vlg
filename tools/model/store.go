package model

import "github.com/timshannon/badgerhold/v4"

const dataDirectory = "data"

// GetStore returns a badgerhold store
func GetStore(readonly bool) (*badgerhold.Store, error) {
	options := badgerhold.DefaultOptions
	options.Dir = dataDirectory
	options.ValueDir = dataDirectory
	options.ReadOnly = readonly
	return badgerhold.Open(options)
}
