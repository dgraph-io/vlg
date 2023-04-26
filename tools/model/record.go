package model

import (
	"fmt"
	"reflect"
	"strings"

	country "github.com/mikekonan/go-countries"
	"github.com/pkg/errors"
	"github.com/timshannon/badgerhold/v4"
)

type Record struct {
	NodeID       string           `csv:"node_id" predicate:"Record.nodeID"`
	Name         string           `csv:"name" predicate:"Record.name"`
	InternalID   string           `csv:"internal_id"  predicate:"Record.internalID"`
	SourceID     string           `csv:"sourceID" predicate:"Record.sourceID"`
	Notes        string           `csv:"notes" predicate:"Record.notes"`
	CountryCodes StringArray      `csv:"country_codes" predicate:"Record.countryCodes"`
	ValidUntil   EmbeddedDateTime `csv:"valid_until" predicate:"Record.validUntil"`
}

type IRecord interface {
}

// Normalize normalizes the record
func (obj *Record) Normalize() {
	if obj.InternalID == "None" {
		obj.InternalID = ""
	}
	if val, ok := country.ByAlpha3Code("XXX"); ok {
		fmt.Println(val)
	}
	for i := range obj.CountryCodes {
		// standardize 'None' country codes for 'XXX'
		if obj.CountryCodes[i] == "None" {
			obj.CountryCodes[i] = "XXX"
			continue
		}
		// if code is not valid, set it to 'XXX'
		val := strings.ToUpper(obj.CountryCodes[i])
		if _, ok := country.ByAlpha3Code(country.Alpha3Code(val)); !ok {
			obj.CountryCodes[i] = "XXX"
		}
	}
	switch {
	case strings.Contains(obj.SourceID, "Panama Papers"):
		obj.SourceID = "PanamaPapers"
	case strings.Contains(obj.SourceID, "Bahamas Leaks"):
		obj.SourceID = "BahamasLeaks"
	case strings.Contains(obj.SourceID, "Paradise Papers"):
		obj.SourceID = "ParadisePapers"
	case strings.Contains(obj.SourceID, "Pandora Papers"):
		obj.SourceID = "PandoraPapers"
	case strings.Contains(obj.SourceID, "Offshore Leaks"):
		obj.SourceID = "OffshoreLeaks"
	default:
		obj.SourceID = "OffshoreLeaks"
	}
}

// RecordByID returns any record by ID. It also returns the type of the record.
func RecordByID(store *badgerhold.Store, id string) (IRecord, string, error) {
	a := []any{
		Entity{},
		Other{},
		Address{},
		Intermediary{},
		Officer{},
	}
	for i := range a {
		elemType := reflect.TypeOf(a[i])
		elem := reflect.New(elemType)
		record := elem.Interface()
		err := store.Get(id, record)
		if err != nil {
			if err == badgerhold.ErrNotFound {
				continue
			}
			return nil, "", err
		}
		switch t := record.(type) {
		case *Entity:
			return record.(*Entity), "entity", nil
		case *Other:
			return record.(*Other), "other", nil
		case *Address:
			return record.(*Address), "address", nil
		case *Intermediary:
			return record.(*Intermediary), "intermediary", nil
		case *Officer:
			return record.(*Officer), "officer", nil
		default:
			return nil, "", errors.Errorf("unknown record type %s", t)
		}
	}
	return nil, "", errors.New("no record found")
}
