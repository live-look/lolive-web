package camforchat

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/authboss"
)

// UserDbStorage represents logic of user storage
type UserDbStorage struct {
	db *sqlx.DB
}

// NewUserDbStorage creates storer object with given db connection
func NewUserDbStorage(db *sqlx.DB) *UserDbStorage {
	return &UserDbStorage{db: db}
}

// New returns empty User object
func (s UserDbStorage) New(ctx context.Context) authboss.User {
	return &User{}
}

// Save Updates user
func (s UserDbStorage) Save(ctx context.Context, user authboss.User) error {
	usr := user.(*User)
	u := &User{}

	findStatement := `SELECT * FROM users WHERE lower(email) = lower($1)`
	err := s.db.Get(u, findStatement, usr.Email)
	if err == sql.ErrNoRows {
		return authboss.ErrUserNotFound
	}
	if err != nil {
		return err
	}

	updateStatement := `UPDATE users
	  SET name = :name,
	  confirmed = :confirmed,
	  password = :password,
	  confirm_selector = :confirm_selector,
	  confirm_verifier = :confirm_verifier,
	  updated_at = NOW()
	  WHERE lower(email) = lower(:email)`

	_, err = s.db.NamedExec(updateStatement,
		map[string]interface{}{
			"name":             usr.Name,
			"email":            usr.Email,
			"confirmed":        usr.Confirmed,
			"password":         usr.Password,
			"confirm_selector": usr.ConfirmSelector,
			"confirm_verifier": usr.ConfirmVerifier,
		})

	return err
}

// Load returns User for given identity (email)
func (s UserDbStorage) Load(ctx context.Context, key string) (authboss.User, error) {
	u := &User{}

	findStatement := `SELECT * FROM users WHERE lower(email) = lower($1)`
	err := s.db.Get(u, findStatement, key)

	if err == sql.ErrNoRows {
		return nil, authboss.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// FindByID finds user by ID
func (s UserDbStorage) FindByID(id int64) (*User, error) {
	u := &User{}

	findStatement := `SELECT * FROM users WHERE id = $1`
	err := s.db.Get(u, findStatement, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Create saves user into database
func (s UserDbStorage) Create(ctx context.Context, user authboss.User) error {
	usr := user.(*User)
	u := &User{}

	findStatement := `SELECT * FROM users WHERE lower(email) = lower($1)`
	err := s.db.Get(u, findStatement, usr.Email)
	if err != sql.ErrNoRows && err != nil {
		return err
	}
	if u.ID != 0 {
		return authboss.ErrUserFound
	}

	// Create user if OK
	insertStatement := `INSERT INTO users (
		name,
		email,
		password,
		confirmed,
		confirm_selector,
		confirm_verifier,
		updated_at,
		created_at) VALUES
		  (:name, lower(:email), :password, :confirmed, :confirm_selector, :confirm_verifier, NOW(), NOW())`

	_, err = s.db.NamedExec(insertStatement,
		map[string]interface{}{
			"name":             usr.Name,
			"email":            usr.Email,
			"confirmed":        usr.Confirmed,
			"password":         usr.Password,
			"confirm_selector": usr.ConfirmSelector,
			"confirm_verifier": usr.ConfirmVerifier,
		})

	return err
}

// LoadByConfirmSelector implements logic of confirmation: loads user by confirm hash
func (s UserDbStorage) LoadByConfirmSelector(ctx context.Context, selector string) (authboss.ConfirmableUser, error) {
	u := &User{}

	findStatement := `SELECT * FROM users WHERE confirm_selector = $1`
	err := s.db.Get(u, findStatement, selector)
	if err == sql.ErrNoRows {
		return nil, authboss.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}
