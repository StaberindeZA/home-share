package user

import (
	"database/sql"
	"log"
)

type LiteLogic struct {
	db *sql.DB
}

func (ll LiteLogic) FindOrCreate(id int, name, email string) (User, error) {
	var u User
	findByIdQuery := `SELECT id, name, email, created_at, updated_at  FROM users WHERE id = ?`
	err := ll.db.QueryRow(findByIdQuery, id).Scan(&u.Id, &u.Name, &u.Email, &u.createdAt, &u.updatedAt)

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

func NewLiteLogic(db *sql.DB) LiteLogic {
	return LiteLogic{
		db: db,
	}
}
