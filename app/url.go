package app

// MapURL represents the mapping of a given
// "long" URL to it's shortened version.
type MapURL struct {
	URLID int64  `db:"url_id"`
	Short string `db:"url_short"`
	Long  string `db:"url_long"`
}

// Long represents the "long" url
// to be shortened
type Long struct {
	URL string `json:"url"`
}
