package camforchat

import (
	"time"
)

// User represent any users
type User struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string

	Confirmed       bool   `db:"confirmed"`
	ConfirmSelector string `db:"confirm_selector"`
	ConfirmVerifier string `db:"confirm_verifier"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GetNamereturns user name
func (u *User) GetName() string {
	return u.Name
}

// GetEmail returns user email
func (u *User) GetEmail() string {
	return u.Email
}

// GetConfirmed returns confirmed flag
func (u *User) GetConfirmed() bool {
	return u.Confirmed
}

// GetConfirmSelector returns confirm selector - hash for confirmation
func (u *User) GetConfirmSelector() string {
	return u.ConfirmSelector
}

// GetConfirmVerifier returns string for verify email
func (u *User) GetConfirmVerifier() string {
	return u.ConfirmVerifier
}

// GetPassword return password string
func (u *User) GetPassword() string {
	return u.Password
}

// PutPassword changes password field
func (u *User) PutPassword(password string) {
	u.Password = password
}

// PutEmail changes email field
func (u *User) PutEmail(email string) {
	u.Email = email
}

// PutConfirmed changes confirmed flag
func (u *User) PutConfirmed(confirmed bool) {
	u.Confirmed = confirmed
}

// PutConfirmSelector changes confirmation selector
func (u *User) PutConfirmSelector(confirmSelector string) {
	u.ConfirmSelector = confirmSelector
}

// PutConfirmVerifier changes confirmation verifier
func (u *User) PutConfirmVerifier(confirmVerifier string) {
	u.ConfirmVerifier = confirmVerifier
}

// GetPID returns unique property of user for identify in system
func (u *User) GetPID() string {
	return u.Email
}

// PutPID changes unique identify (email for us)
func (u *User) PutPID(email string) {
	u.Email = email
}

// GetArbitrary returns map of additional user fields such as Name
func (u *User) GetArbitrary() map[string]string {
	return map[string]string{
		"name": u.Name,
	}
}

// PutArbitrary changes additional user fields
func (u *User) PutArbitrary(arbitrary map[string]string) {
	if n, ok := arbitrary["name"]; ok {
		u.Name = n
	}
}
