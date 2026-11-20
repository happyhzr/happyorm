package happyorm

import (
	"database/sql"
	"fmt"

	"github.com/happyhzr/happyorm/dialect"
	"github.com/happyhzr/happyorm/log"
	"github.com/happyhzr/happyorm/session"
)

type TxFunc func(*session.Session) (any, error)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver string, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s not found", driver)
		return nil, fmt.Errorf("dialect %s not found", driver)
	}
	e := &Engine{db: db, dialect: dial}
	log.Info("connect database success")
	return e, err
}

func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Error("failed to close database")
		return
	}
	log.Info("close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

func (e *Engine) Transaction(f TxFunc) (result any, err error) {
	s := e.NewSession()
	errBegin := s.Begin()
	if errBegin != nil {
		return nil, errBegin
	}
	defer func() {
		r := recover()
		if r != nil {
			s.Rollback()
			panic(r)
		} else if err != nil {
			s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}
