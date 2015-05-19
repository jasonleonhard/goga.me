package dal

import (
	_ "github.com/lib/pq"
	"testing"
)

func newMatchForTest(t *testing.T) *Match {
	return NewMatch(newDbForTest(t))
}

func TestMatchCRUD(t *testing.T) {
	black, _ := newUserRowForTest(t)
	white, _ := newUserRowForTest(t)
	m := newMatchForTest(t)
	matchRow, err := m.BeginMatch(nil, black, white)
	if err != nil {
		t.Errorf("Beginning a match with two users should work. Error: %v", err)
	}

	if matchRow.BlackUserID <= 0 || matchRow.WhiteUserID <= 0 {
		t.Fatal("Users must have positive ids")
	}

}

func TestMatchOneUser(t *testing.T) {
	black, _ := newUserRowForTest(t)
	m := newMatchForTest(t)
	matchRow, err := m.BeginMatch(nil, black, black)
	if err != nil {
		t.Errorf("Beginning a match with one user should work. Error: %v", err)
	}

	if matchRow.BlackUserID <= 0 || matchRow.WhiteUserID <= 0 {
		t.Fatal("Users must have positive ids")
	}
}
