package model

type Relationship struct {
	FromID           int64    `csv:"node_id_start"`
	FromType         string   `csv:"-"`
	ToID             int64    `csv:"node_id_end"`
	ToType           string   `csv:"-"`
	RelationshipType string   `csv:"rel_type"`
	Status           string   `csv:"status"`
	StartDate        DateTime `csv:"start_date"`
	EndDate          DateTime `csv:"end_date"`
	SourceID         string   `csv:"source_id"`
}
