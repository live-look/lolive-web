package camforchat

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewDb(t *testing.T) {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	_, err := NewDb(spec)
	assert.Nil(t, err)
}

func TestGetDb(t *testing.T) {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := NewDb(spec)

	assert.Nil(t, err)

	ctx := context.WithValue(context.Background(), ctxKeyDb, db)

	c, ok := GetDb(ctx)
	assert.True(t, ok)
	assert.IsType(t, c, db)
}
