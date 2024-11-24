package models

type Source struct {
	Name  string `json:"name"`
	Types []Type
} //	@name	Source

//nolint:gochecknoglobals //ok
var Sources = []Source{
	ManualSource,
	SteamSource,
	GoodreadsSource,
}

//nolint:gochecknoglobals //ok
var ManualSource = Source{
	Name: "Manual",
	Types: []Type{
		AmountType,
	},
}

//nolint:gochecknoglobals //ok
var SteamSource = Source{
	Name: "Steam",
	Types: []Type{
		SteamCompletionPercentage,
		CompletedGames,
	},
}

//nolint:gochecknoglobals //ok
var GoodreadsSource = Source{
	Name: "Goodreads",
	Types: []Type{
		FinishedBooks,
	},
}
