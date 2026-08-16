package entry

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DummyLogic struct{}

func (dl DummyLogic) Create(userId int, start, end time.Time) (string, error) {
	fmt.Println(userId, start, end)

	return "ok", nil
}

func (dl DummyLogic) Read(id int) (Entry, error) {
	fmt.Println("Read id:", id)
	return Entry{
		id:     id,
		userId: 999,
		value:  Talking,
		start:  time.Now(),
		end:    time.Now(),
	}, nil
}

func (dl DummyLogic) Update(id int, value EntryValue) (string, error) {
	fmt.Println("Update id:", id, value)

	return "ok", nil
}

func (dl DummyLogic) Delete(id int) (string, error) {
	fmt.Println("Delete id:", id)

	return "ok", nil
}

func (dl DummyLogic) List(userId int, start, end time.Time) ([]Entry, error) {
	fmt.Println("List userId", userId)
	entries := make([]Entry, 2)
	entries[0] = Entry{
		id:     123,
		userId: userId,
		start:  time.Now(),
		end:    time.Now(),
	}
	entries[1] = Entry{
		id:     321,
		userId: userId,
		start:  time.Now(),
		end:    time.Now(),
	}

	return entries, nil
}

func NewDummyLogic() DummyLogic {
	return DummyLogic{}
}

type MVPLogic struct {
	db *sql.DB
}

func (ml MVPLogic) Create(userId int, start, end time.Time) (string, error) {
	insertEntriesQuery := `INSERT INTO entries (user_id, value, start, end) VALUES (?, ?, ?, ?);`
	if _, err := ml.db.Exec(insertEntriesQuery, userId, Talking, start, end); err != nil {
		return "error", err
	}

	return "ok", nil
}

func (ml MVPLogic) Read(id int) (Entry, error) {
	var e Entry
	readQuery := `SELECT id, user_id, value, start, end FROM entries WHERE id = ?;`
	err := ml.db.QueryRow(readQuery, id).Scan(&e.id, &e.userId, &e.value, &e.start, &e.end)

	return e, err
}

func (ml MVPLogic) Update(id int, value EntryValue) (string, error) {
	updateQuery := `UPDATE entries SET value = ?, updated_at = ? WHERE id = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := ml.db.ExecContext(ctx, updateQuery, value, time.Now(), id)
	if err != nil {
		log.Fatalf("failed to execute update: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("failed to look up rows affected: %v", err)
	}

	if rowsAffected != 1 {
		log.Fatalf("unexpected number of rows affected: %d", rowsAffected)
	}

	return "ok", nil
}

func (ml MVPLogic) Delete(id int) (string, error) {
	query := `DELETE FROM entries WHERE id = ?`
	if _, err := ml.db.Exec(query, id); err != nil {
		return "error", err
	}

	return "ok", nil
}

func (ml MVPLogic) List(userId int, start, end time.Time) ([]Entry, error) {
	query := `SELECT id, user_id, value, start, end FROM entries WHERE user_id = ? and start <= ? and end >= ?`
	rows, err := ml.db.Query(query, userId, end, start)
	if err != nil {
		return []Entry{}, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.id, &e.userId, &e.value, &e.start, &e.end); err != nil {
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
