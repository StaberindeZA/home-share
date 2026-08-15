package user

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type LiteLogic struct {
	db *sql.DB
}

func (ll LiteLogic) FindOrCreate(id int, name, email string) (User, error) {
	var u User
	findByIDQuery := `SELECT id, name, email, created_at, updated_at  FROM users WHERE id = ?`
	err := ll.db.QueryRow(findByIDQuery, id).Scan(&u.Id, &u.Name, &u.Email, &u.createdAt, &u.updatedAt)

	if err == sql.ErrNoRows {
		findByEmailQuery := `SELECT id, name, email, created_at, updated_at  FROM users WHERE email = ?`
		err = ll.db.QueryRow(findByEmailQuery, email).Scan(&u.Id, &u.Name, &u.Email, &u.createdAt, &u.updatedAt)
	}

	if err == sql.ErrNoRows {
		insertUsersQuery := `INSERT INTO users (name, email) VALUES (?, ?) RETURNING id, name, email, created_at, updated_at;`

		err := ll.db.QueryRow(insertUsersQuery, name, email).Scan(&u.Id, &u.Name, &u.Email, &u.createdAt, &u.updatedAt)
		if err != nil {
			log.Fatalf("Error inserting new record: %v", err)
		}
	}

	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Query failed: %v", err)
	}

	return u, nil
}

func (ll LiteLogic) FindByEmail(email string) (User, error) {
	var u User
	findByEmailQuery := `SELECT id, name, email, created_at, updated_at  FROM users WHERE email = ?`
	if err := ll.db.QueryRow(findByEmailQuery, email).Scan(&u.Id, &u.Name, &u.Email, &u.createdAt, &u.updatedAt); err != nil {
		return User{}, err
	}

	return u, nil
}

func (ll LiteLogic) UpdateByEmail(email, name string) (string, error) {
	updateQuery := `UPDATE users SET name = ?, updated_at = ? WHERE email = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := ll.db.ExecContext(ctx, updateQuery, name, time.Now(), email)
	if err != nil {
		log.Printf("failed to execute update: %v", err)
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to look up rows affected: %v", err)
		return "", err
	}

	if rowsAffected != 1 {
		log.Printf("unexpected number of rows affected: %d", rowsAffected)
		return "", errors.New("unexpected number of rows affected")
	}

	return "ok", nil
}

func NewLiteLogic(db *sql.DB) LiteLogic {
	return LiteLogic{
		db: db,
	}
}
