package query

import (
	"github.com/jmoiron/sqlx"
)

// Executor is a wrapper capable of running raw sql commands,
// start transactions, etc.
type Executor struct {
	noCopy
	sqlx.Ext
}

// NewExecutor returns a new instance of a SQL executor,
// based on the provided connection.
func NewExecutor(db *sqlx.DB) *Executor {
	return &Executor{Ext: db}
}

// NewBuilder returns a new query builder, preconfigured to use
// this executor.
func (e *Executor) NewBuilder() *Builder {
	return &Builder{executor: e}
}
