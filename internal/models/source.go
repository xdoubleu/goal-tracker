package models

type Source struct {
	Name  string `json:"name"`
	Types []Type
} //	@name	Source

//nolint:gochecknoglobals //ok
var Sources = []Source{
	SteamSource,
	GoodreadsSource,
}

//nolint:gochecknoglobals //ok
var SteamSource = Source{
	Name: "Steam",
	Types: []Type{
		SteamCompletionRate,
	},
}

//nolint:gochecknoglobals //ok
var GoodreadsSource = Source{
	Name: "Goodreads",
	Types: []Type{
		FinishedBooks,
	},
}
