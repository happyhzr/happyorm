package happyorm

import (
	"database/sql"

	"github.com/happyhzr/happyorm/log"
	"github.com/happyhzr/happyorm/session"
)

type Engine struct {
	db *sql.DB
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
	e := &Engine{db: db}
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
	return session.New(e.db)
}
