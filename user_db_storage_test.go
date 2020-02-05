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

func TestUserDbStorageLoad(t *testing.T) {
	spec := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := NewDb(spec)
	assert.Nil(t, err)

	s := NewUserDbStorage(db)
	_, err = s.Load(context.Background(), "thisuserdoesnotexist@example.com")

	assert.Equal(t, authboss.ErrUserNotFound, err)

	user := &User{
		Name:            "foo",
		Email:           "foo3@example.com",
		Confirmed:       true,
		Password:        "qwerty",
		ConfirmSelector: "fizzbuzz",
		ConfirmVerifier: "aaa",
	}

	err = s.Create(context.Background(), user)
	assert.Nil(t, err)

	usr, err := s.Load(context.Background(), "foo3@example.com")
	assert.Nil(t, err)

	user = usr.(*User)

	assert.Equal(t, "foo", user.GetName())
	assert.Equal(t, "foo3@example.com", user.GetEmail())
	assert.True(t, user.GetConfirmed())
}
