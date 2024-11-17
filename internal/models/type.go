package models

type Type struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
} //	@name	Type

var AmountType = Type{
	ID:   0,
	Name: "Amount",
}

var SteamCompletionPercentage = Type{
	ID:   1,
	Name: "Steam completion percentage",
}

var ActualCompletionPercentage = Type{
	ID:   2,
	Name: "Actual completion percentage",
}

var CompletedGames = Type{
	ID:   3,
	Name: "Completed games",
}

var FinishedBooks = Type{
	ID:   4,
	Name: "Finished books",
}
