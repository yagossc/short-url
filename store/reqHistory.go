package store

import (
	"github.com/yagossc/short-url/app"
	"github.com/yagossc/short-url/query"
)

// InsertReq saves a new ReqHistory entry in the database.
func InsertReq(db *query.Executor, req *app.ReqHistory) error {
	sql := db.NewBuilder()

	sql.Add("INSERT INTO req_history(url_short, req_time)")
	sql.Add("VALUES ($1, $2)")
	sql.SetParam(1, req.ShortURL)
	sql.SetParam(2, req.ReqTime)

	_, err := sql.Exec()
	if err != nil {
		return err
	}

	return nil
}

// FindAllReqByShort retrieves all requests in history given a short URL
func FindAllReqByShort(db *query.Executor, short string) ([]app.ReqHistory, error) {
	sql := db.NewBuilder()

	sql.Add(" SELECT *")
	sql.From("  FROM req_history")
	sql.Where("WHERE url_short = ?", short)

	var result []app.ReqHistory
	err := sql.Select(&result)

	return result, err
}
