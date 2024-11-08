package models

type Source struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
} //	@name	Source

var Sources = []Source{
	ManualSource,
}

var ManualSource = Source{
	ID:   0,
	Name: "Manual",
}
