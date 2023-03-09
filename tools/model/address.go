package model

type Location struct {
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
	Source string  `json:"source"`
}

type Address struct {
	Record

	Address string `csv:"address"`

	Location *Location `json:"location"`
}
