package model

import "strings"

// The gocarina/gocsv package supports the declaration of UnmarshalCSV methods.
type StringArray []string

// UnmarshalCSV unmarshals a string array from a semicolon-delimited string.
func (array *StringArray) UnmarshalCSV(csv string) (err error) {
	*array = strings.Split(csv, ";")
	return nil
}
