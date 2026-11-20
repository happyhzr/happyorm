package happyorm

import (
	"database/sql"
	"fmt"
	"strings"

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

func difference(a []string, b []string) []string {
	diff := []string{}
	mapB := map[string]bool{}
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		_, ok := mapB[v]
		if !ok {
			diff = append(diff, v)
		}
	}
	return diff
}

func (e *Engine) Migrate(value any) error {
	_, err := e.Transaction(func(s *session.Session) (any, error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s does not exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		rows.Close()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			_, err := s.Raw(sqlStr).Exec()
			if err != nil {
				return nil, err
			}
		}
		if len(delCols) == 0 {
			return nil, nil
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err := s.Exec()
		return nil, err
	})
	return err
}
