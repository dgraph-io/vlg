package model

type Intermediary struct {
	Record

	Status  string `csv:"status"`
	Address string `csv:"address"`
}
