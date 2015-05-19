package dal

import (
	_ "github.com/lib/pq"
	"testing"
)

func newUserForTest(t *testing.T) *User {
	return NewUser(newDbForTest(t))
}

func newUserRowForTest(t *testing.T) (*UserRow, error) {
	u := newUserForTest(t)
	userRow, err := u.Signup(nil, newEmailForTest(), "abc123", "abc123")
	return userRow, err
}

func TestUserCRUD(t *testing.T) {
	userRow, err := newUserRowForTest(t)

	// Signup
	if err != nil {
		t.Errorf("Signing up user should work. Error: %v", err)
	}
	if userRow == nil {
		t.Fatal("Signing up user should work.")
	}
	if userRow.ID <= 0 {
		t.Fatal("Signing up user should work.")
	}
	if userRow.CreatedAt == nil {
		t.Fatal("Created at should be non-nil")
	}
	if userRow.UpdatedAt == nil {
		t.Fatal("Updated at should be non-nil")
	}

}

func TestUserDelete(t *testing.T) {
	u := newUserForTest(t)
	userRow, err := u.Signup(nil, newEmailForTest(), "abc123", "abc123")

	// DELETE FROM users WHERE id=...
	_, err = u.DeleteById(nil, userRow.ID)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}
