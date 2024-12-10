package entities

type Song struct {
	ID          int     `json:"id" db:"id"`
	Group       string  `json:"group" db:"group"`
	Song        string  `json:"song" db:"song"`
	ReleaseDate *string `json:"releaseDate" db:"release_date"`
	Link        *string `json:"link" db:"link"`
	Text        *string `json:"text" db:"text"`
}
