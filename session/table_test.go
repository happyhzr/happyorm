package session

import (
	"log"
	"testing"
)

type User struct {
	Name string `happyorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	err := s.DropTable()
	log.Println(err)
	err = s.CreateTable()
	log.Println(err)
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
