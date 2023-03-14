package user

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// data ...
type data struct {
	XUID       string `db:"xuid"`
	Username   string `db:"username"`
	Role       string `db:"role"`
	Kills      int    `db:"kills"`
	Deaths     int    `db:"deaths"`
	CPS        bool   `db:"cps"`
	ScoreBoard bool   `db:"scoreboard"`
}

// SelectData returns the *data of the player in the given database.
func (u *User) SelectData(db *sqlx.DB) *data {
	data := &data{
		CPS:        true,
		ScoreBoard: true,
	}
	rows, err := db.Queryx("SELECT * FROM playerdata WHERE xuid = ?", u.XUID())
	defer func() {
		go func() {
			if rows != nil {
				rows.Close()
			}
		}()
	}()
	if err != nil {
		u.Disconnect(err)
	}
	for rows.Next() {
		_ = rows.StructScan(&data)
	}

	return data
}

// UpdateData updates the data of the player to the given database
func (u *User) UpdateData(db *sqlx.DB) {
	_, err := db.Exec("REPLACE INTO playerdata VALUES ($1, $2, $3, $4, $5, $6, $7)",
		u.XUID(),
		u.Name(),
		u.rank.Name(),
		u.data.CPS,
		u.data.ScoreBoard,
		u.data.Kills,
		u.data.Deaths)
	if err != nil {
		fmt.Println(err)
	}
}
