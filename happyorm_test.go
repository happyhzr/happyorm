package happyorm

import (
	"errors"
	"testing"

	"github.com/happyhzr/happyorm/session"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "happy.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `happyorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (any, error) {
		s.Model(&User{}).CreateTable()
		s.Insert(&User{"Tom", 18})
		return nil, errors.New("error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (any, error) {
		s.Model(&User{}).CreateTable()
		_, errInsert := s.Insert(&User{"Tom", 18})
		return nil, errInsert
	})
	u := &User{}
	s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
