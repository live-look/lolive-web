package models

import (
	"database/sql"
	"fmt"
	txdb "github.com/DATA-DOG/go-txdb"
	"github.com/romanyx/polluter"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func prepareBroadcastTestDb(t *testing.T) (db *sql.DB, closeConn func() error) {
	dbSpec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("CAMFORCHAT_APP_POSTGRES_USER"),
		os.Getenv("CAMFORCHAT_APP_POSTGRES_PASSWORD"),
		os.Getenv("CAMFORCHAT_APP_POSTGRES_HOST"),
		"camforchat")

	txdb.Register("psql_txdb", "postgres", dbSpec)

	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())

	db, err := sql.Open("psql_txdb", cName)
	if err != nil {
		t.Fatalf("open psql_txdb connection: %s", err)
	}

	seed, err := os.Open("../fixtures/users.yml")
	if err != nil {
		t.Fatalf("failed to open seed file: %s", err)
	}
	defer seed.Close()

	p := polluter.New(polluter.PostgresEngine(db))
	if err := p.Pollute(seed); err != nil {
		t.Fatalf("failed to pollute: %s", err)
	}

	return db, db.Close
}

func TestCreate(t *testing.T) {
}
