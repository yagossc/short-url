package query

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// ErrTooManyRows is a signal error indicating that a resultset
// returned more rows than expected.
var ErrTooManyRows = errors.New("too many rows")

// ErrNoRows is a signal error indication that a resultset did not
// return any rows.
var ErrNoRows = errors.New("no rows returned")

// Builder is a auxiliary buffer to more easily build and
// execute SQL queries
type Builder struct {
	executor *Executor

	selectSQL strings.Builder
	fromSQL   strings.Builder
	whereSQL  strings.Builder

	params []interface{}
}

// Add unconditionally appends an sql string into the builder's buffer.
// If any values are supplyed, the '?' parameters will be replaced by the
// appropriate token.
func (b *Builder) Add(sql string, values ...interface{}) {
	s := b.loadParameters(sql, values)
	b.appendSQL(&b.selectSQL, s)
}

// From appends the sql string in a special buffer, used to mount the
// 'from' clause. Otherwise, it has the same mechanics as the Add method.
func (b *Builder) From(sql string, values ...interface{}) {
	s := b.loadParameters(sql, values)
	b.appendSQL(&b.fromSQL, s)
}

// Where appends the sql string in a special buffer, used to mount the
// 'where' clause. Otherwise, it has the same mechanics as the Add method.
func (b *Builder) Where(sql string, values ...interface{}) {
	s := b.loadParameters(sql, values)
	b.appendSQL(&b.whereSQL, s)
}

// One executes the query, mapping the first and only row
// to the specified interface. If the query returns zero or
// more than one row, an error is returned.
func (b *Builder) One(value interface{}) error {
	sql := b.buildSQL()

	rows, err := b.executor.Queryx(sql, b.params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			return ErrTooManyRows
		}

		v := reflect.ValueOf(value).Elem()
		if v.Type().Kind() == reflect.Struct {
			if err := rows.StructScan(value); err != nil {
				return err
			}
		} else {
			if err := rows.Scan(value); err != nil {
				return err
			}
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if count == 0 {
		return ErrNoRows
	}

	return nil
}

func (b *Builder) appendSQL(sb *strings.Builder, s string) {
	if s == "" {
		return
	}

	sb.WriteString(s)
	sb.WriteString(" ")
}

func (b *Builder) buildSQL() string {
	return b.selectSQL.String() + b.fromSQL.String() + b.whereSQL.String()
}

func (b *Builder) loadParameters(originalSQL string, values []interface{}) string {
	newSQL := originalSQL
	psize := len(b.params)

	for _, value := range values {
		b.params = append(b.params, value)
		psize++
		newSQL = strings.Replace(newSQL, "?", "$"+strconv.Itoa(psize), 1)
	}

	return newSQL
}
