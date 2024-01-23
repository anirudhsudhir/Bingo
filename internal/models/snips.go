package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snip struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnipModel struct {
	DB *sql.DB
}

func (m *SnipModel) InsertSnip(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snips (title, content, created, expires)
	VALUES (?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnipModel) ReadSnip(id int) (*Snip, error) {
	stmt := `SELECT id, title, content, created, expires FROM snips
	WHERE id = ? AND expires > UTC_TIMESTAMP()`
	row := m.DB.QueryRow(stmt, id)

	snip := &Snip{}
	err := row.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Content, &snip.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return snip, nil
}

func (m *SnipModel) GetLatestSnips() ([]*Snip, error) {
	stmt := `SELECT id, title, content, created, expires FROM snips
	WHERE expires > UTC_TIMESTAMP ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snips := []*Snip{}
	for rows.Next() {
		snip := &Snip{}
		err := rows.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)
		if err != nil {
			return nil, err
		}
		snips = append(snips, snip)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snips, nil
}
