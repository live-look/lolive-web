package models

import (
	"github.com/jmoiron/sqlx"
)

type broadcastDatabaseStorage struct {
	db *sqlx.DB
}

func NewBroadcastDbStorage(db *sqlx.DB) *broadcastStorer {
	return &broadcastDatabaseStorage{db: db}
}

func (bs *broadcastDatabaseStorage) Save(broadcast *Broadcast) error {
	insertQuery := `INSERT INTO broadcasts (id, user_id, state, created_at) VALUES ($1, $2)`
	if err := bs.db.Exec(
		insertQuery,
		broadcast.ID,
		broadcast.UserID,
		broadcast.State,
		broadcast.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

// UpdateState changes state of broadcast
func (bs *broadcastDatabaseStorage) UpdateState(state BroadcastState) error {
	updateQuery := `UPDATE broadcasts SET state = :state WHERE id = :id`
	_, err := bs.db.NamedExec(updateQuery,
		map[string]interface{}{
			"state": state,
			"id":    b.ID,
		})

	return err
}

// FindByState retrive list of online broadcasts
func (bs *broadcastDatabaseStorage) FindByState(state BroadcastState) ([]*Broadcast, error) {
	var broadcasts []*Broadcast
	selectQuery := `SELECT b.*, u.name AS user_name FROM broadcasts b INNER JOIN users u ON u.id = b.user_id WHERE state = $1 ORDER BY created_at DESC`

	err := bs.db.Select(&broadcasts, selectQuery, state)

	return broadcasts, err
}

// Find gets broadcast from db by ID
func (bs *broadcastDatabaseStorage) Find(ID string) (*Broadcast, error) {
	broadcast := &Broadcast{}

	selectQuery := `SELECT b.*, u.name AS user_name
					FROM broadcasts b
					INNER JOIN users u
					  ON u.id = b.user_id
					WHERE b.id = $1`

	err := bs.db.Get(broadcast, selectQuery, ID)

	return broadcast, err
}
