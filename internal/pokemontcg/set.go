package pokemontcg

type Set struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Series      string `json:"series"`
	Total       int    `json:"total"`
	ReleaseDate string `json:"releaseDate"`
}
