package models

type Type struct {
	ID       int64  `json:"id"`
	SourceID int64  `json:"sourceId"`
	Name     string `json:"name"`
} //	@name	Type

var Types = []Type{
	AmountType,
	HoursType,
}

var AmountType = Type{
	ID:       0,
	SourceID: ManualSource.ID,
	Name:     "Amount",
}

var HoursType = Type{
	ID:       1,
	SourceID: ManualSource.ID,
	Name:     "Hours",
}
