package main

import (
	"log"

	"github.com/happyhzr/happyorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := happyorm.NewEngine("sqlite3", "happy.db")
	defer engine.Close()
	s := engine.NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) VALUES (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	log.Printf("exec success, %d affected\n", count)
}
