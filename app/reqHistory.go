package app

// ReqHistory defines the relation between a
// short URL and its requests as a time series
type ReqHistory struct {
	ReqID    int64  `db:"req_id"`
	ShortURL string `db:"url_short"`
	ReqTime  int64  `db:"req_time"`
}

// Short represents the "short" URL
// to be queried in the history
type Short struct {
	URL string `json:"url"`
}
