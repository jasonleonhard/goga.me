package dal

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewMatch(db *sqlx.DB) *Match {
	match := &Match{}
	match.db = db
	match.table = "matches"
	match.hasID = true
	return match
}

type MatchRow struct {
	ID          int64      `db:"id"`
	CreatedAt   *time.Time `db:"created_at"`
	BlackUserID int64      `db:"black_user_id"`
	WhiteUserID int64      `db:"white_user_id"`
}

type Match struct {
	Base
}

// BeginMatch takes two users, black and white, and then returns a MatchRow
func (m *Match) BeginMatch(tx *sqlx.Tx, black, white *UserRow) (*MatchRow, error) {
	data := make(map[string]interface{})
	data["black_user_id"] = black.ID
	if white != nil {
		data["white_user_id"] = white.ID
	}

	sqlResult, err := m.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return m.matchRowFromSqlResult(tx, sqlResult)
}

func (m *Match) GetById(tx *sqlx.Tx, id int64) (*MatchRow, error) {
	match := &MatchRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", m.table)
	err := m.db.Get(match, query, id)

	return match, err
}

func (m *Match) matchRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*MatchRow, error) {
	matchId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.GetById(tx, matchId)
}
