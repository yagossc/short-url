package app

// MapURL represents the mapping of a given
// "long" URL to it's shortened version.
type MapURL struct {
	URLID    int64  `db:"url_id"`
	ShortURL string `db:"url_short"`
	URL      string `db:"url_long"`
}

// LongURL represents the "long" url
// to be shortened
type LongURL struct {
	URL string `json:"url"`
}
