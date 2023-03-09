package model

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/timshannon/badgerhold/v4"
)

type Record struct {
	NodeID       int64            `csv:"node_id"`
	InternalID   string           `csv:"internal_id"`
	SourceID     int64            `csv:"source_id"`
	Notes        string           `csv:"notes"`
	Name         string           `csv:"name"`
	CountryCodes StringArray      `csv:"country_codes"`
	ValidUntil   EmbeddedDateTime `csv:"valid_until"`
}

type DateTime struct {
	time.Time
}

var badDateMap = map[string]string{
	"01-APR-1001": "01-APR-2001",
	"04-AUG-0095": "04-AUG-1995",
	"06-DEC-1194": "06-DEC-1994",
	"09-AUG-0004": "09-AUG-2004",
	"11-JAN-0990": "11-JAN-1999",
	"12-MAR-0096": "12-MAR-1996",
	"18-DEC-2812": "18-DEC-2012",
	"18-JUN-0082": "18-JUN-1982",
	"24-JUN-0003": "24-JUN-2003",
	"30-DEC-0203": "30-DEC-2003",
	"30-OCT-0014": "30-OCT-2014",
	"18/19/2015":  "18-DEC-2015",
	"29/02/2006":  "28-FEB-2006", // 29th Feb 2006 doesn't exist
	"31/02/2013":  "28-FEB-2013", // 31st Feb 2013 doesn't exist
	"06/01.2009":  "06-JAN-2009", // WTF?
	"26/08-2010":  "26-AUG-2010", // Have you no shame?
	"03/03/3008":  "03-MAR-2008", // 3008 doesn't exist
	"28.02/2014":  "28-FEB-2014", // Sigh
	"28/02/201":   "28-FEB-2011",
	"16/20/2013":  "01-JAN-2013",
	"1212012":     "12-DEC-2012",
	"01.004.04":   "01-JAN-2004",
	"08/11/204":   "11-AUG-2004",
	"02/11/5005":  "11-FEB-2005",
	"284/2015":    "01-JAN-2015",
	"14-MAR-2998": "14-MAR-1998",
	"22-MAY-5012": "22-MAY-2012",
	"30-MAR-3012": "30-MAR-2012",
	"06-AUG-3003": "06-AUG-2003",
	"24-MAR-1200": "24-MAR-2000",
	"10-JAN-0992": "10-JAN-1992",
	"28-AUG-5015": "28-AUG-2015",
	"29-JAN-0104": "29-JAN-2004",
	"03-JUN-1203": "03-JUN-2013",
	"23-AUG-1013": "23-AUG-2013",
	"30-DEC-1759": "30-DEC-1959",
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return nil
	}
	// handle a variety of "bad" dates
	csv = strings.Replace(csv, "0004", "2004", 1)
	csv = strings.Replace(csv, "0006", "2006", 1)
	csv = strings.Replace(csv, "0009", "2009", 1)
	csv = strings.Replace(csv, "0011", "2011", 1)
	csv = strings.Replace(csv, "0012", "2012", 1)
	csv = strings.Replace(csv, "0015", "2015", 1)
	csv = strings.Replace(csv, "0199", "1999", 1)
	csv = strings.Replace(csv, "0200", "2000", 1)
	csv = strings.Replace(csv, "0201", "2001", 1)
	csv = strings.Replace(csv, "0202", "2002", 1)
	csv = strings.Replace(csv, "0206", "2006", 1)
	csv = strings.Replace(csv, "0208", "2008", 1)
	csv = strings.Replace(csv, "0213", "2013", 1)
	csv = strings.Replace(csv, "0214", "2014", 1)
	csv = strings.Replace(csv, "1196", "1996", 1)
	csv = strings.Replace(csv, "1198", "1998", 1)
	if replaceDate, ok := badDateMap[csv]; ok {
		csv = replaceDate
	}
	// handle a variety of date formats
	formats := []string{
		"02-Jan-2006",
		"Jan 02, 2006",
		"2006-01-02",
		"01/02/2006",
		"1/2/2006",
		"01/02/06",
		"01.02.2006",
		"01.2.2006",
		"1.02.2006",
		"2006",
		"06/01/02",
		"02/01/06",
		"02/01/2006",
		"01-02-2006",
		"2/1/06",
		"2/01/2006",
		"02012006",
		"01022006",
		"1/2/06",
		"02-012006",
		"2-1-06",
		"Jan-2006",
	}
	// Title case the month name
	csv = strings.Title(strings.ToLower(csv))
	for _, format := range formats {
		date.Time, err = time.Parse(format, csv)
		if err == nil {
			if date.Time.Year() < 1800 || date.Time.Year() > 2800 {
				return errors.Errorf("invalid date %s", csv)
			}
			return nil
		}
	}
	return errors.Wrapf(err, "failed to parse date %s", csv)
}

type EmbeddedDateTime struct {
	time.Time
}

// Convert the CSV string as internal date
func (date *EmbeddedDateTime) UnmarshalCSV(csv string) (err error) {
	value := regexp.MustCompile(`\d{4}`).FindString(csv)
	if value == "" {
		return nil
	}
	date.Time, err = time.Parse("2006", value)
	return err
}

type StringArray []string

func (array *StringArray) UnmarshalCSV(csv string) (err error) {
	*array = strings.Split(csv, ";")
	return nil
}

func RecordByID(store *badgerhold.Store, id int64) (any, error) {
	a := []any{
		Entity{},
		Other{},
		Address{},
		Intermediary{},
		Officer{},
	}
	for i := range a {
		elemType := reflect.TypeOf(a[i])
		elemSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
		x := reflect.New(elemSlice.Type())
		x.Elem().Set(elemSlice)
		err := store.Find(x.Interface(), badgerhold.Where("NodeID").Eq(id))
		if err != nil {
			return nil, err
		}
		switch x.Elem().Len() {
		case 0:
			continue
		case 1:
			return x.Elem().Index(0).Interface(), nil
		default:
			return nil, errors.New("multiple records found")
		}
	}
	return nil, errors.New("no record found")
}
