package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

type Record struct {
	Timestamp uint64 `json:"unixtime"`
	Posts     uint64 `json:"posts"`
}

type ServerAnswer struct {
	Board   string   `json:"board"`
	Records []Record `json:"records"`
}

func (r *Database) MigrateStats() error {
	query := `
	CREATE TABLE IF NOT EXISTS stats(
		board VARCHAR(255) NOT NULL,
		unixtime INTEGER NOT NULL PRIMARY KEY,
		posts INTEGER DEFAULT 0
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *Database) GetBoardStats(board string) (*ServerAnswer, error) {
	query := `SELECT * FROM stats WHERE board = "%s"`
	query = fmt.Sprintf(query, board)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	answer := ServerAnswer{}
	for rows.Next() {
		var record Record
		err = rows.Scan(
			&answer.Board,
			&record.Timestamp,
			&record.Posts,
		)
		if err != nil {
			return nil, err
		}
		answer.Records = append(answer.Records, record)
	}
	if len(answer.Records) == 0 {
		return nil, errors.New("no records for such board")
	}

	return &answer, nil
}

func (d *Database) InsertRecord(board string, r *Record) error {
	query := `
	INSERT INTO stats(board, unixtime, posts)
	values(?,?,?)
	`
	_, err := d.db.Exec(
		query,
		board,
		r.Timestamp,
		r.Posts,
	)

	return err
}
