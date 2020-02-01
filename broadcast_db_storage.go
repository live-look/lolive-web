package camforchat

import (
	"github.com/jmoiron/sqlx"
)

// BroadcastDbStorage implementation of BroadcastStorer
type BroadcastDbStorage struct {
	db *sqlx.DB
}

// NewBroadcastDbStorage creates new instance of BroadcastDbStorage
func NewBroadcastDbStorage(db *sqlx.DB) *BroadcastDbStorage {
	return &BroadcastDbStorage{db: db}
}

// Save inserts new broadcast row into database
func (bs *BroadcastDbStorage) Save(broadcast *Broadcast) error {
	_, err := bs.db.NamedExec(
		`INSERT INTO broadcasts (id, user_id, state, created_at) VALUES (:id, :user_id, :state, :created_at)`,
		map[string]interface{}{
			"id":         broadcast.ID,
			"user_id":    broadcast.UserID,
			"state":      broadcast.State,
			"created_at": broadcast.CreatedAt,
		},
	)
	return err
}

// UpdateState changes state of broadcast
func (bs *BroadcastDbStorage) UpdateState(broadcast *Broadcast, state BroadcastState) error {
	updateQuery := `UPDATE broadcasts SET state = :state WHERE id = :id`
	_, err := bs.db.NamedExec(updateQuery,
		map[string]interface{}{
			"state": state,
			"id":    broadcast.ID,
		})

	return err
}

// FindByState retrive list of online broadcasts
func (bs *BroadcastDbStorage) FindByState(state BroadcastState) ([]*Broadcast, error) {
	var broadcasts []*Broadcast
	selectQuery := `SELECT b.*, u.name AS user_name FROM broadcasts b INNER JOIN users u ON u.id = b.user_id WHERE state = $1 ORDER BY created_at DESC`

	err := bs.db.Select(&broadcasts, selectQuery, state)

	return broadcasts, err
}

// Find gets broadcast from db by ID
func (bs *BroadcastDbStorage) Find(ID string) (*Broadcast, error) {
	broadcast := &Broadcast{}

	selectQuery := `SELECT b.*, u.name AS user_name
					FROM broadcasts b
					INNER JOIN users u
					  ON u.id = b.user_id
					WHERE b.id = $1`

	err := bs.db.Get(broadcast, selectQuery, ID)

	return broadcast, err
}
