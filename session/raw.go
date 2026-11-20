package session

import (
	"database/sql"
	"strings"

	"github.com/happyhzr/happyorm/clause"
	"github.com/happyhzr/happyorm/dialect"
	"github.com/happyhzr/happyorm/log"
	"github.com/happyhzr/happyorm/schema"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	tx       *sql.Tx
	refTable *schema.Schema
	clause   *clause.Clause
	sql      *strings.Builder
	sqlVars  []any
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect, clause: &clause.Clause{}, sql: &strings.Builder{}}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = &clause.Clause{}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, values ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	result, err := s.DB().Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	rows, err := s.DB().Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return rows, nil
}

type CommonDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)
