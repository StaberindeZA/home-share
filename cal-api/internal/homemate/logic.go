package homemate

import "database/sql"

type LiteLogic struct {
	db *sql.DB
}

func NewLiteLogic(db *sql.DB) LiteLogic {
	return LiteLogic{
		db: db,
	}
}

func (l LiteLogic) Create(homeID, mateID int, role Role) error {
	insertQuery := `INSERT INTO homemates (home_id, mate_id, role) VALUES (?, ?, ?);`
	_, err := l.db.Exec(insertQuery, homeID, mateID, role)
	return err
}

func (l LiteLogic) ReadForHomeAndMate(homeID, mateID int) (HomeMate, error) {
	var hm HomeMate
	query := `SELECT id, home_id, mate_id, role, created_at, updated_at FROM homemates where home_id = ? and mate_id = ?`

	err := l.db.QueryRow(query, homeID, mateID).Scan(&hm.id, &hm.homeID, &hm.mateID, &hm.Role, &hm.createdAt, &hm.updatedAt)

	return hm, err
}

func (l LiteLogic) Delete(homeID, mateID int) error {
	query := `DELETE FROM homemates WHERE home_id = ? AND mate_id = ?`
	_, err := l.db.Exec(query, homeID, mateID)
	return err
}
