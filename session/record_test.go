package session

import "testing"

var user1 = &User{"Tom", 18}
var user2 = &User{"Sam", 25}
var user3 = &User{"Jack", 25}

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	assertUsers := func(t *testing.T, users []User) {
		t.Helper()
		if len(users) != 2 {
			t.Fatalf("Find returned %d users, want 2", len(users))
		}

		wantTom, wantSam := *user1, *user2
		if (users[0] != wantTom || users[1] != wantSam) &&
			(users[0] != wantSam || users[1] != wantTom) {
			t.Fatalf("Find returned %+v, want %+v and %+v", users, wantTom, wantSam)
		}
	}

	t.Run("value_slice", func(t *testing.T) {
		s := testRecordInit(t)
		var users []User
		if err := s.Find(&users); err != nil {
			t.Fatalf("Find failed: %v", err)
		}
		assertUsers(t, users)
	})

	t.Run("pointer_slice", func(t *testing.T) {
		s := testRecordInit(t)
		var users []*User
		if err := s.Find(&users); err != nil {
			t.Fatalf("Find failed: %v", err)
		}
		if len(users) != 2 {
			t.Fatalf("Find returned %d users, want 2", len(users))
		}
		if users[0] == nil || users[1] == nil || users[0] == users[1] {
			t.Fatalf("Find returned %v, want two distinct non-nil pointers", users)
		}
		assertUsers(t, []User{*users[0], *users[1]})
	})
}
