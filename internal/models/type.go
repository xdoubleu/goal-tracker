package models

type Type struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
} //	@name	Type

//nolint:gochecknoglobals //ok
var SteamCompletionRate = Type{
	ID:   0,
	Name: "Steam completion rate",
}

//nolint:gochecknoglobals //ok
var FinishedBooks = Type{
	//nolint:mnd //ok
	ID:   1,
	Name: "Finished books",
}
