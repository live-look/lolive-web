package camforchat

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/authboss"
	"os"
	"testing"
)

func TestNewUserDbStorage(t *testing.T) {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := NewDb(spec)
	assert.Nil(t, err)

	s := NewUserDbStorage(db)

	assert.IsType(t, db, s.db)
}

func TestUserDbStorageNew(t *testing.T) {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := NewDb(spec)
	assert.Nil(t, err)

	s := NewUserDbStorage(db)
	assert.NotNil(t, s.New(context.Background()))
}

func TestUserDbStorageCreate(t *testing.T) {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := NewDb(spec)
	assert.Nil(t, err)

	user := &User{
		Name:            "foo",
		Email:           "foo@example.com",
		Confirmed:       true,
		Password:        "qwerty",
		ConfirmSelector: "fizzbuzz",
		ConfirmVerifier: "aaa",
	}
	s := NewUserDbStorage(db)

	err = s.Create(context.Background(), user)
	assert.Nil(t, err)

	user2 := &User{
		Name:            "bar",
		Email:           "foo@example.com",
		Confirmed:       false,
		Password:        "toor",
		ConfirmSelector: "xxxx",
		ConfirmVerifier: "zzzz",
	}
	err = s.Create(context.Background(), user2)
	assert.Equal(t, authboss.ErrUserFound, err)
}

func TestUserDbStorageSave(t *testing.T) {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := NewDb(spec)
	assert.Nil(t, err)

	user := &User{
		Name:            "foo",
		Email:           "foo2@example.com",
		Confirmed:       true,
		Password:        "qwerty",
		ConfirmSelector: "fizzbuzz",
		ConfirmVerifier: "aaa",
	}
	s := NewUserDbStorage(db)

	assert.Equal(t, authboss.ErrUserNotFound, s.Save(context.Background(), user))
}
