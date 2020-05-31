package store

import (
	"github.com/yagossc/short-url/app"
	"github.com/yagossc/short-url/query"
)

// FindURLByShort returns an URL instance with specified SHORT field.
func FindURLByShort(db *query.Executor, short string) (*app.MapURL, error) {

	sql := db.NewBuilder()
	sql.Add("SELECT  url_id,")
	sql.Add("        url_short,")
	sql.Add("        url_long")
	sql.Add("   FROM url")
	sql.Where("WHERE url_short = ?", short)

	var url app.MapURL
	if err := sql.One(&url); err == query.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &url, nil
}
