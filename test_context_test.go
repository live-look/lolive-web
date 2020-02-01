package camforchat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func setupTestContext(t *testing.T) *testContext {
	dbSpec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := sqlx.Connect("pgx", dbSpec)
	assert.Nil(t, err)

	err = db.Ping()
	assert.Nil(t, err)

	webrtc := NewWebrtc()

	return newTestContext(context.Background(), db, webrtc)
}

func userFixture(t *testing.T, ctx *testContext, email string) *User {
	userStorer := NewUserStorer(ctx.db)
	usr := &User{Name: email, Email: email, Password: "qwerty1234567890"}
	err := userStorer.Create(ctx.context, usr)
	assert.Nil(t, err)

	u, err := userStorer.Load(ctx.context, email)
	assert.Nil(t, err)

	return u.(*User)
}

func broadcastFixture(user *User) *Broadcast {
	return &Broadcast{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		UserName:  user.Name,
		State:     BroadcastStateOffline,
		CreatedAt: time.Now(),
	}
}
