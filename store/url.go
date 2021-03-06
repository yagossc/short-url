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
	sql.From("  FROM url_map")
	sql.Where("WHERE url_short = ?", short)

	var url app.MapURL
	if err := sql.One(&url); err == query.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &url, nil
}

// FindAllURL retrieves all available url mappings
func FindAllURL(db *query.Executor) ([]app.MapURL, error) {
	sql := db.NewBuilder()

	sql.Add("SELECT *")
	sql.From(" FROM url_map")

	var result []app.MapURL
	err := sql.Select(&result)
	return result, err
}

// InsertURL saves a new MapURL in the database.
func InsertURL(db *query.Executor, url *app.MapURL) (string, error) {
	sql := db.NewBuilder()

	sql.Add("INSERT INTO url_map(url_short, url_long)")
	sql.Add("VALUES ($1, $2)")
	sql.SetParam(1, url.Short)
	sql.SetParam(2, url.Long)

	_, err := sql.Exec()
	if err != nil {
		return "", err
	}

	return url.Short, nil
}
