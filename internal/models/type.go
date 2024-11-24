package models

type Type struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
} //	@name	Type

//nolint:gochecknoglobals //ok
var AmountType = Type{
	ID:   0,
	Name: "Amount",
}

//nolint:gochecknoglobals //ok
var SteamCompletionPercentage = Type{
	ID:   1,
	Name: "Steam completion percentage",
}

//nolint:gochecknoglobals //ok
var CompletedGames = Type{
	//nolint:mnd //ok
	ID:   2,
	Name: "Completed games",
}

//nolint:gochecknoglobals //ok
var FinishedBooks = Type{
	//nolint:mnd //ok
	ID:   3,
	Name: "Finished books",
}
