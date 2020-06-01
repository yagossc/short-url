package dbtest

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // needed to load sqlite3 driver
	"github.com/yagossc/short-url/query"
)

// SQLTest defines a way to hold multiple
// sequences of queries for testing
type SQLTest map[string][]string

// GenDefaultLoad generates a db load for testing
func GenDefaultLoad() []string {
	return []string{
		`CREATE TABLE url_map (
		   url_id INTEGER NOT NULL,
		   url_short TEXT NOT NULL,
		   url_long TEXT  NOT NULL,

		   CONSTRAINT pk_url PRIMARY KEY (url_id),
		   CONSTRAINT uq_url_short UNIQUE (url_short)
		 );`,
		`CREATE TABLE req_history (
		   req_id INTEGER NOT NULL,
		   url_short TEXT NOT NULL,
		   req_time BIGINT NOT NULL,

		   CONSTRAINT pk_req PRIMARY KEY (req_id),
		   CONSTRAINT fk_url FOREIGN KEY (url_short) REFERENCES url_map (url_short)
		 );`,
		`INSERT INTO url_map(url_short, url_long)
		 VALUES ('yxZ8byjhRui', 'http://www.google.com')`,
		`INSERT INTO req_history(url_short, req_time)
		 VALUES ('yxZ8byjhRui', '1590977472');`,
		`INSERT INTO req_history(url_short, req_time)
		 VALUES ('yxZ8byjhRui', '1590977473');`,
		`INSERT INTO req_history(url_short, req_time)
		 VALUES ('yxZ8byjhRui', '1590977474');`,
		`SELECT * FROM req_history`,
	}
}

// DBHandler is a handler function that accepts a db instance
type DBHandler func(db *query.Executor)

// WithDB runs the specified handler providing a valid sqlite3 database.
func WithDB(handler DBHandler) {
	// Please see https://www.sqlite.org/inmemorydb.html
	//
	// Look for the section "In-memory Databases And Shared Cache" to understand why
	// the url cannot be a simple ":memory:"
	db, err := sqlx.Open("sqlite3", "file:memdb?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	handler(query.NewExecutor(db))
}

// Mock seeds a db with predefined values for testing
func Mock(db *query.Executor, sqls []string) error {
	for _, query := range sqls {
		i, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("query numer %d, while running command '%s'\nerror: %v", i, query, err)
		}
	}

	return nil
}
