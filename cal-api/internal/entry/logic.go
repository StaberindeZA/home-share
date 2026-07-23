package entry

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type DummyLogic struct{}

func (dl DummyLogic) Create(userId string, start, end time.Time) (string, error) {
	fmt.Println(userId, start, end)

	if userId == "fail" {
		return "", errors.New("things failed")
	}

	return "ok", nil
}

func (dl DummyLogic) Read(id string) (Entry, error) {
	fmt.Println("Read id:", id)
	return Entry{
		id:     id,
		userId: "123_userid",
		start:  time.Now(),
		end:    time.Now(),
	}, nil
}

func (dl DummyLogic) Delete(id string) (string, error) {
	fmt.Println("Delete id:", id)

	return "ok", nil
}

func (dl DummyLogic) List(userId string, start, end time.Time) ([]Entry, error) {
	fmt.Println("List userId", userId)
	entries := make([]Entry, 2)
	entries[0] = Entry{
		id:     "123",
		userId: userId,
		start:  time.Now(),
		end:    time.Now(),
	}
	entries[1] = Entry{
		id:     "321",
		userId: userId,
		start:  time.Now(),
		end:    time.Now(),
	}

	return entries, nil
}

func NewDummyLogic() DummyLogic {
	return DummyLogic{}
}

type MVPLogic struct{
	db *sql.DB
}

func (ml MVPLogic) Create(userId string, start, end time.Time) (string, error) {
	insertEntriesQuery := `INSERT INTO entries (user_id, start, end) VALUES (?, ?, ?);`
	if _, err := ml.db.Exec(insertEntriesQuery, userId, start, end); err != nil {
		return "error", err
	}

	return "ok", nil
}

func (ml MVPLogic) Read(id string) (Entry, error) {
	fmt.Println("Read id:", id)
	return Entry{
		id:     id,
		userId: "123_userid",
		start:  time.Now(),
		end:    time.Now(),
	}, nil
}

func (ml MVPLogic) Delete(id string) (string, error) {
	query := `DELETE FROM entries WHERE id = ?`
	if _, err := ml.db.Exec(query, id); err != nil {
		return "error", err
	}

	return "ok", nil
}

func (ml MVPLogic) List(userId string, start, end time.Time) ([]Entry, error) {
	query := `SELECT id, user_id, start, end FROM entries WHERE user_id = ? and start <= ? and end >= ?`
	rows, err := ml.db.Query(query, userId, end, start)
	if err != nil {
		return []Entry{}, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.id, &e.userId, &e.start, &e.end); err != nil {
			log.Fatal(err)
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func NewMVPLogic(db *sql.DB) MVPLogic {
	return MVPLogic{
		db: db,
	}
}
