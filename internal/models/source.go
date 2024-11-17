package models

type Source struct {
	Name  string `json:"name"`
	Types []Type
} //	@name	Source

var Sources = []Source{
	ManualSource,
	SteamSource,
	GoodreadsSource,
}

var ManualSource = Source{
	Name: "Manual",
	Types: []Type{
		AmountType,
	},
}

var SteamSource = Source{
	Name: "Steam",
	Types: []Type{
		SteamCompletionPercentage,
		ActualCompletionPercentage,
		CompletedGames,
	},
}

var GoodreadsSource = Source{
	Name: "Goodreads",
	Types: []Type{
		FinishedBooks,
	},
}
