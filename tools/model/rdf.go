package model

import (
	"bytes"
	"fmt"
	"io"

	tstore "github.com/matthewmcneely/triplestore"
)

func RDFEscapeLiteral(l string) string {
	var buf bytes.Buffer
	for _, r := range l {
		switch r {
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		default:
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func RDFLiteralEmpty(literal tstore.Literal) bool {
	switch literal.Type() {
	case tstore.XsdString:
		switch literal.Value() {
		case "":
			return true
		case "0001-01-01 00:00:00 +0000 UTC":
			return true
		default:
			return false
		}
	case tstore.XsdDateTime:
		return literal.Value() == "0001-01-01 00:00:00 +0000 UTC"
	default:
		return false
	}
}

func RDFEncodeTriples(w io.Writer, triples []tstore.Triple) {
	for _, tri := range triples {
		literal, _ := tri.Object().Literal()
		if !RDFLiteralEmpty(literal) {
			fmt.Fprintf(w, "%s <%s> \"%s\"^^<%s> .\n", tri.Subject(), tri.Predicate(), RDFEscapeLiteral(literal.Value()), literal.Type())
		}
	}
}
