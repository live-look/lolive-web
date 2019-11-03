package models

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func prepareTestViewerDb(t *testing.T) (db *sqlx.DB, closeConn func() error) {
	dbSpec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := sqlx.Connect("postgres", dbSpec)
	if err != nil {
		t.Fatalf("open postgresql connection: %s", err)
	}

	return db, db.Close
}

func TestViewerCreate(t *testing.T) {
	db, closeConn := prepareTestViewerDb(t)
	defer closeConn()

	userStorer := NewUserStorer(db)
	usr := &User{Name: "test", Email: "foo@example.com", Password: "xyz"}
	err := userStorer.Create(context.Background(), usr)
	if err != nil {
		t.Errorf("Error creating user: %s", err)
		return
	}

	u, err := userStorer.Load(context.Background(), "foo@example.com")
	if err != nil {
		t.Errorf("Error loading user: %s", err)
		return
	}

	user := u.(*User)
	b := NewBroadcast(db, user.ID)
	err = b.Create()
	if err != nil {
		t.Errorf("Error creating broadcast: %s", err)
		return
	}

	v := NewViewer(db, user.ID, b.ID)
	err = v.Create()
	if err != nil {
		t.Errorf("Error creating viewer: %s", err)
		return
	}
}
