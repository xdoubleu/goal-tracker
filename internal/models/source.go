package models

type Source struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Types []Type
} //	@name	Source

var Sources = []Source{
	ManualSource,
	SteamSource,
	GoodreadsSource,
}

var ManualSource = Source{
	ID:   0,
	Name: "Manual",
	Types: []Type{
		AmountType,
	},
}

var SteamSource = Source{
	ID:   1,
	Name: "Steam",
	Types: []Type{
		SteamCompletionPercentage,
		ActualCompletionPercentage,
		CompletedGames,
	},
}

var GoodreadsSource = Source{
	ID:   2,
	Name: "Goodreads",
	Types: []Type{
		FinishedBooks,
	},
}
