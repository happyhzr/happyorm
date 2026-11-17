package happyorm

import (
	"database/sql"
	"fmt"

	"github.com/happyhzr/happyorm/dialect"
	"github.com/happyhzr/happyorm/log"
	"github.com/happyhzr/happyorm/session"
)

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
