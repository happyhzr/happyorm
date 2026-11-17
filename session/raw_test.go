package session

import (
	"database/sql"
	"os"
	"testing"

	"github.com/happyhzr/happyorm/dialect"

	_ "github.com/mattn/go-sqlite3"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../happy.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}
