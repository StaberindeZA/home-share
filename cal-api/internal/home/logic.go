package home

import (
	"database/sql"
	"log"

	"cal-api/internal/homemate"
	"cal-api/internal/user"

	"github.com/gosimple/slug"
)

type LiteLogic struct {
	db        *sql.DB
	userLogic user.UserLogic
}

func NewLiteLogic(db *sql.DB, userLogic user.UserLogic) LiteLogic {
	return LiteLogic{
		db:        db,
		userLogic: userLogic,
	}
}

func (l LiteLogic) Create(name, description string) (string, error) {
	homeSlug := slug.Make(name)

	insertQuery := `INSERT INTO homes (name, slug, description) VALUES (?, ?, ?);`
	_, err := l.db.Exec(insertQuery, name, homeSlug, description)

	return homeSlug, err
}

func (l LiteLogic) Read(slug string) (Home, error) {
	var h Home
	readQuery := `SELECT id, name, slug, description, created_at, updated_at FROM homes WHERE slug = ?;`
	err := l.db.QueryRow(readQuery, slug).Scan(&h.ID, &h.name, &h.slug, &h.description, &h.createdAt, &h.updatedAt)

	return h, err
}

func (l LiteLogic) Delete(slug string) error {
	query := `DELETE FROM homes WHERE slug = ?`
	_, err := l.db.Exec(query, slug)
	return err
}

func (l LiteLogic) List(mateID int) ([]Home, error) {
	query := `SELECT h.id, h.name, h.slug, h.description, h.created_at, h.updated_at
	FROM homes as h
	JOIN homemates hm on h.id = hm.home_id
	where hm.mate_id = ?`

	rows, err := l.db.Query(query, mateID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var homes []Home
	var home Home

	for rows.Next() {
		// Scan directly matching your SELECT order: h.id, hm.id, u.email, u.name
		err := rows.Scan(&home.ID, &home.name, &home.slug, &home.description, &home.createdAt, &home.updatedAt)
		if err != nil {
			log.Fatal(err)
		}

		homes = append(homes, home)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		var emptyHomes []Home
		return emptyHomes, err
	}

	return homes, nil
}

func (l LiteLogic) ReadMates(slug string) ([]ListHomeMatesMates, error) {
	query := `SELECT u.id, u.email, u.name, hm.role
	FROM homes as h
	JOIN homemates hm on h.id = hm.home_id
	JOIN users u on hm.mate_id = u.id
	where h.slug = ?`

	rows, err := l.db.Query(query, slug)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var homeMatesMates []ListHomeMatesMates
	//firstRow := true

	for rows.Next() {
		var id int
		var email, name string
		var role homemate.Role

		// Scan directly matching your SELECT order: h.id, hm.id, u.email, u.name
		err := rows.Scan(&id, &email, &name, &role)
		if err != nil {
			log.Fatal(err)
		}

		// Create and append the expanded reference
		mate := ListHomeMatesMates{
			ID:    id,
			Email: email,
			Name:  name,
			Role:  role,
		}
		homeMatesMates = append(homeMatesMates, mate)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return homeMatesMates, err
	}

	return homeMatesMates, nil
}
