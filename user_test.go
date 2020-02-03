package camforchat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserGetEmail(t *testing.T) {
	usr := &User{Email: "foo@example.com"}
	assert.Equal(t, "foo@example.com", usr.GetEmail())
}

func TestUserGetConfirmed(t *testing.T) {
	usr := &User{Confirmed: true}
	assert.True(t, usr.GetConfirmed())
}

func TestUserGetConfirmSelector(t *testing.T) {
	usr := &User{ConfirmSelector: "foobar"}
	assert.Equal(t, "foobar", usr.GetConfirmSelector())
}

func TestUserGetConfirmVerifier(t *testing.T) {
	usr := &User{ConfirmVerifier: "foobar"}
	assert.Equal(t, "foobar", usr.GetConfirmVerifier())
}

func TestUserGetPassword(t *testing.T) {
	usr := &User{Password: "foobar"}
	assert.Equal(t, "foobar", usr.GetPassword())
}

func TestUserPutPassword(t *testing.T) {
	usr := &User{Password: "foobar"}

	usr.PutPassword("fizzbuzz")
	assert.Equal(t, "fizzbuzz", usr.GetPassword())
}

func TestUserPutEmail(t *testing.T) {
	usr := &User{Email: "foo@example.com"}

	usr.PutEmail("fizzbuzz@example.com")
	assert.Equal(t, "fizzbuzz@example.com", usr.GetEmail())
}

func TestUserPutConfirmed(t *testing.T) {
	usr := &User{Confirmed: true}

	usr.PutConfirmed(false)
	assert.False(t, usr.GetConfirmed())
}

func TestUserPutConfirmSelector(t *testing.T) {
	usr := &User{ConfirmSelector: "foobar"}

	usr.PutConfirmSelector("fizzbuzz")
	assert.Equal(t, "fizzbuzz", usr.GetConfirmSelector())
}

func TestUserPutConfirmVerifier(t *testing.T) {
	usr := &User{ConfirmVerifier: "foobar"}

	usr.PutConfirmVerifier("fizzbuzz")
	assert.Equal(t, "fizzbuzz", usr.GetConfirmVerifier())
}

func TesUserGetPID(t *testing.T) {
	usr := &User{Email: "foobar@example.com"}

	assert.Equal(t, "foobar@example.com", usr.GetPID)
}

func TesUserPutPID(t *testing.T) {
	usr := &User{Email: "foobar@example.com"}

	usr.PutPID("fizzbuzz@example.com")
	assert.Equal(t, "fizzbuzz@example.com", usr.GetPID)
}

func TestUserGetArbitrary(t *testing.T) {
	usr := &User{Name: "foobar"}

	assert.Equal(t, usr.GetArbitrary()["name"], "foobar")
}

func TestUserPutArbitrary(t *testing.T) {
	usr := &User{Name: "foobar"}

	arbitrary := make(map[string]string)
	arbitrary["name"] = "fizzbuzz"

	usr.PutArbitrary(arbitrary)

	assert.Equal(t, usr.GetArbitrary()["name"], "fizzbuzz")
}
